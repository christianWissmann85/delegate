#!/bin/bash

# create-bundle.sh - Create markdown bundles from directories
# Usage: ./create-bundle.sh <output-file> <dir1> [dir2] [dir3] ...

set -e

# Check arguments
if [ $# -lt 2 ]; then
    echo "Usage: $0 <output-file> <dir1> [dir2] [dir3] ..."
    echo "Example: $0 tmp/bundles/internal-bundle.md internal/"
    exit 1
fi

OUTPUT_FILE="$1"
shift  # Remove first argument, leaving only directories

# Create output directory if it doesn't exist
OUTPUT_DIR=$(dirname "$OUTPUT_FILE")
mkdir -p "$OUTPUT_DIR"

# Initialize the bundle file
echo "# Code Bundle" > "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "Generated on: $(date)" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Function to get file size in KB
get_size_kb() {
    local file="$1"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        stat -f%z "$file" | awk '{print int($1/1024)}'
    else
        stat -c%s "$file" | awk '{print int($1/1024)}'
    fi
}

# Function to generate tree structure
generate_tree() {
    local dir="$1"
    local prefix="$2"
    local items=($(ls -1 "$dir" 2>/dev/null | sort))
    local count=${#items[@]}
    
    for ((i=0; i<$count; i++)); do
        local item="${items[$i]}"
        local path="$dir/$item"
        
        if [ $i -eq $((count-1)) ]; then
            echo "${prefix}└── $item"
            local new_prefix="${prefix}    "
        else
            echo "${prefix}├── $item"
            local new_prefix="${prefix}│   "
        fi
        
        if [ -d "$path" ]; then
            generate_tree "$path" "$new_prefix"
        fi
    done
}

# Generate file tree
echo "## File Tree" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo '```' >> "$OUTPUT_FILE"
for dir in "$@"; do
    if [ -d "$dir" ]; then
        echo "$(basename "$dir")/"
        generate_tree "$dir" ""
    fi
done
echo '```' >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Function to categorize Go files
categorize_go_file() {
    local file="$1"
    
    # Check if it's a test file
    if [[ "$file" == *_test.go ]]; then
        echo "test"
        return
    fi
    
    # Check file content for interfaces
    if grep -q "^type.*interface {" "$file" 2>/dev/null; then
        echo "interface"
        return
    fi
    
    # Check for main package
    if grep -q "^package main" "$file" 2>/dev/null; then
        echo "main"
        return
    fi
    
    # Default to implementation
    echo "implementation"
}

# Collect all files by type
declare -a interface_files
declare -a type_files
declare -a implementation_files
declare -a test_files
declare -a other_files

for dir in "$@"; do
    if [ -d "$dir" ]; then
        while IFS= read -r -d '' file; do
            if [[ "$file" == *.go ]]; then
                category=$(categorize_go_file "$file")
                case "$category" in
                    "interface")
                        interface_files+=("$file")
                        ;;
                    "test")
                        test_files+=("$file")
                        ;;
                    *)
                        implementation_files+=("$file")
                        ;;
                esac
            else
                other_files+=("$file")
            fi
        done < <(find "$dir" -type f -print0)
    fi
done

# Function to append file to bundle
append_file() {
    local file="$1"
    local relative_path="${file#./}"
    
    echo "## $relative_path" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    
    # Determine language for syntax highlighting
    local ext="${file##*.}"
    local lang=""
    case "$ext" in
        go) lang="go" ;;
        js|jsx) lang="javascript" ;;
        ts|tsx) lang="typescript" ;;
        py) lang="python" ;;
        sh) lang="bash" ;;
        yml|yaml) lang="yaml" ;;
        json) lang="json" ;;
        md) lang="markdown" ;;
        *) lang="" ;;
    esac
    
    echo '```'"$lang" >> "$OUTPUT_FILE"
    cat "$file" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    echo '```' >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
}

# Add files in order: interfaces, types, implementations, tests, others
echo "## Code Files" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [ ${#interface_files[@]} -gt 0 ]; then
    echo "### Interfaces" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    for file in "${interface_files[@]}"; do
        append_file "$file"
    done
fi

if [ ${#implementation_files[@]} -gt 0 ]; then
    echo "### Implementations" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    for file in "${implementation_files[@]}"; do
        append_file "$file"
    done
fi

if [ ${#test_files[@]} -gt 0 ]; then
    echo "### Tests" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    for file in "${test_files[@]}"; do
        append_file "$file"
    done
fi

if [ ${#other_files[@]} -gt 0 ]; then
    echo "### Other Files" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    for file in "${other_files[@]}"; do
        append_file "$file"
    done
fi

# Check file size and warn if too large
TOTAL_SIZE=$(get_size_kb "$OUTPUT_FILE")
echo "Bundle created: $OUTPUT_FILE (${TOTAL_SIZE}KB)"

if [ $TOTAL_SIZE -gt 500 ]; then
    echo "WARNING: Bundle exceeds 500KB! Consider splitting into smaller bundles."
fi