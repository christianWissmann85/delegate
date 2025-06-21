#!/bin/bash

# Specialized bundle.md creator for the delegate repository
# Creates organized documentation with table of contents, file tree, and contents

set -e

# Output file
OUTPUT="bundle.md"

# Color codes for terminal output (optional)
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[*]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Function to check if file should be included
should_include() {
    local file="$1"
    
    # Skip the output file itself
    [[ "$file" == "$OUTPUT" ]] && return 1
    
    # Skip executables (no extension files that are likely binaries)
    [[ "$file" == "delegate" ]] && return 1
    [[ "$file" == "e2e" ]] && return 1
    
    # Skip .env file explicitly
    [[ "$file" == ".env" ]] && return 1
    
    # Skip .delegate/outputs directory contents
    [[ "$file" =~ ^\.delegate/outputs/ ]] && return 1
    
    # Check gitignore
    if git check-ignore "$file" &> /dev/null; then
        return 1
    fi
    
    return 0
}

# Function to categorize files for better organization
get_file_category() {
    local file="$1"
    
    if [[ "$file" =~ ^docs/ ]]; then
        echo "Documentation"
    elif [[ "$file" =~ ^internal/.*\.go$ ]]; then
        echo "Go Source - Internal"
    elif [[ "$file" =~ ^internal/.*_test\.go$ ]]; then
        echo "Go Tests - Internal"
    elif [[ "$file" =~ ^testdata/ ]]; then
        echo "Test Data"
    elif [[ "$file" == "main.go" ]]; then
        echo "Go Source - Main"
    elif [[ "$file" =~ \.(md|MD)$ ]]; then
        echo "Documentation"
    elif [[ "$file" =~ ^\.claude/ ]]; then
        echo "Configuration - Claude"
    elif [[ "$file" == "go.mod" ]] || [[ "$file" == "go.sum" ]]; then
        echo "Go Module Files"
    elif [[ "$file" =~ \.(sh|bash)$ ]]; then
        echo "Scripts"
    elif [[ "$file" == ".gitignore" ]] || [[ "$file" == "LICENSE" ]]; then
        echo "Project Meta"
    else
        echo "Other"
    fi
}

# Function to generate enhanced table of contents
generate_toc() {
    echo "## Table of Contents"
    echo ""
    echo "### Quick Navigation"
    echo "- [Project Overview](#project-overview)"
    echo "- [File Tree](#file-tree)"
    echo "- [Documentation Files](#documentation-files)"
    echo "  - [Project Documentation](#project-documentation)"
    echo "  - [Architecture Documentation](#architecture-documentation)"
    echo "  - [Development Documentation](#development-documentation)"
    echo "  - [Guides](#guides)"
    echo "  - [Reference Documentation](#reference-documentation)"
    echo "- [Source Code](#source-code)"
    echo "  - [Main Application](#main-application)"
    echo "  - [Internal Packages](#internal-packages)"
    echo "  - [Configuration](#configuration-code)"
    echo "  - [Tests](#tests)"
    echo "- [Configuration Files](#configuration-files)"
    echo "- [Scripts](#scripts)"
    echo ""
}

print_status "Creating specialized bundle.md for delegate repository..."

# Start creating the bundle file
{
    echo "# Delegate Repository Bundle"
    echo ""
    echo "Complete source code and documentation bundle for the Delegate project."
    echo ""
    echo "**Generated on:** $(date '+%Y-%m-%d %H:%M:%S %Z')"
    echo "**Repository:** delegate"
    echo "**Type:** Go MCP (Model Context Protocol) Server"
    echo ""
    
    # Add project overview section
    echo "## Project Overview"
    echo ""
    echo "Delegate is a Model Context Protocol (MCP) server that enables Large Language Models (LLMs) to interact with other LLMs through a standardized interface."
    echo ""
    
    # Generate table of contents
    generate_toc
    
    # File tree section
    echo "## File Tree"
    echo ""
    echo '```'
    tree -a -I '.git|.env|delegate|e2e' --gitignore
    echo '```'
    echo ""
    
    # Documentation section
    echo "## Documentation Files"
    echo ""
    
    # Project documentation
    echo "### Project Documentation"
    echo ""
    
    # Main project docs
    for file in README.md CLAUDE.md LICENSE; do
        if [[ -f "$file" ]] && should_include "$file"; then
            echo "#### $file"
            echo ""
            echo '```markdown'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Architecture documentation
    echo "### Architecture Documentation"
    echo ""
    
    for file in docs/architecture/*.md; do
        if [[ -f "$file" ]] && should_include "$file"; then
            echo "#### $file"
            echo ""
            echo '```markdown'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Development documentation
    echo "### Development Documentation"
    echo ""
    
    for file in docs/development/*.md; do
        if [[ -f "$file" ]] && should_include "$file"; then
            echo "#### $file"
            echo ""
            echo '```markdown'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Guides
    echo "### Guides"
    echo ""
    
    for file in docs/guides/*.md; do
        if [[ -f "$file" ]] && should_include "$file"; then
            echo "#### $file"
            echo ""
            echo '```markdown'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Reference documentation
    echo "### Reference Documentation"
    echo ""
    
    for file in docs/reference/*.md; do
        if [[ -f "$file" ]] && should_include "$file"; then
            echo "#### $file"
            echo ""
            echo '```markdown'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Source code section
    echo "## Source Code"
    echo ""
    
    # Main application
    echo "### Main Application"
    echo ""
    
    if [[ -f "main.go" ]]; then
        echo "#### main.go"
        echo ""
        echo '```go'
        cat "main.go"
        echo '```'
        echo ""
    fi
    
    # Go module files
    echo "### Go Module Files"
    echo ""
    
    for file in go.mod go.sum; do
        if [[ -f "$file" ]]; then
            echo "#### $file"
            echo ""
            echo '```'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Internal packages
    echo "### Internal Packages"
    echo ""
    
    # Process internal packages in order
    for package in config extractor handlers logger mcp models providers storage; do
        package_dir="internal/$package"
        if [[ -d "$package_dir" ]]; then
            echo "#### Package: $package"
            echo ""
            
            # List Go files in this package (excluding tests)
            find "$package_dir" -name "*.go" -not -name "*_test.go" -type f | sort | while read -r file; do
                if should_include "$file"; then
                    echo "##### $file"
                    echo ""
                    echo '```go'
                    cat "$file"
                    echo '```'
                    echo ""
                fi
            done
            
            # Handle subdirectories (like providers/anthropic, providers/google, etc.)
            find "$package_dir" -mindepth 1 -type d | sort | while read -r subdir; do
                subpackage=$(basename "$subdir")
                find "$subdir" -name "*.go" -not -name "*_test.go" -type f | sort | while read -r file; do
                    if should_include "$file"; then
                        echo "##### $file"
                        echo ""
                        echo '```go'
                        cat "$file"
                        echo '```'
                        echo ""
                    fi
                done
            done
        fi
    done
    
    # Tests section
    echo "### Tests"
    echo ""
    
    # Find all test files
    find internal -name "*_test.go" -type f | sort | while read -r file; do
        if should_include "$file"; then
            echo "#### $file"
            echo ""
            echo '```go'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Configuration files
    echo "## Configuration Files"
    echo ""
    
    # .gitignore
    if [[ -f ".gitignore" ]]; then
        echo "### .gitignore"
        echo ""
        echo '```'
        cat ".gitignore"
        echo '```'
        echo ""
    fi
    
    # Claude settings
    if [[ -f ".claude/settings.local.json" ]]; then
        echo "### .claude/settings.local.json"
        echo ""
        echo '```json'
        cat ".claude/settings.local.json"
        echo '```'
        echo ""
    fi
    
    # Scripts section
    echo "## Scripts"
    echo ""
    
    # Find all shell scripts
    find . -name "*.sh" -type f | sort | while read -r file; do
        if should_include "$file"; then
            echo "### $file"
            echo ""
            echo '```bash'
            cat "$file"
            echo '```'
            echo ""
        fi
    done
    
    # Test data section (if any files exist)
    if [[ -d "testdata" ]] && [[ -n "$(ls -A testdata 2>/dev/null)" ]]; then
        echo "## Test Data"
        echo ""
        
        find testdata -type f | sort | while read -r file; do
            if should_include "$file"; then
                echo "### $file"
                echo ""
                # Detect file type for syntax highlighting
                case "${file##*.}" in
                    json) lang="json" ;;
                    xml) lang="xml" ;;
                    yaml|yml) lang="yaml" ;;
                    *) lang="" ;;
                esac
                echo '```'"$lang"
                cat "$file"
                echo '```'
                echo ""
            fi
        done
    fi
    
} > "$OUTPUT"

# Summary statistics
total_files=$(grep -c '^###\+ ' "$OUTPUT" || echo "0")
file_size=$(du -h "$OUTPUT" | cut -f1)
go_files=$(find . -name "*.go" -type f | wc -l)
test_files=$(find . -name "*_test.go" -type f | wc -l)
doc_files=$(find . -name "*.md" -type f | wc -l)

print_success "Successfully created $OUTPUT"
echo ""
echo "📊 Bundle Statistics:"
echo "   📄 Total files included: $total_files"
echo "   💾 Bundle file size: $file_size"
echo "   🔵 Go source files: $go_files"
echo "   🧪 Test files: $test_files"
echo "   📚 Documentation files: $doc_files"
echo ""
print_status "Bundle includes all source code and documentation"
print_warning "Excluded: .env, binary files (delegate, e2e), .delegate/outputs/"