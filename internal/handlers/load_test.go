package handlers_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/christianwissmann85/delegate/internal/extractor"
	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/storage"
)

// TestConcurrentInvokeCalls tests concurrent invoke operations
func TestConcurrentInvokeCalls(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	store, err := storage.NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	extractFactory := extractor.NewFactory()
	providerFactory := &mockProviderFactory{}
	invokeHandler := handlers.NewInvokeHandler(providerFactory, store, extractFactory)

	ctx := context.Background()

	// Number of concurrent calls
	numCalls := 20
	var wg sync.WaitGroup
	errors := make(chan error, numCalls)
	outputIDs := make(chan string, numCalls)

	// Track timing
	start := time.Now()

	// Launch concurrent invokes
	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req := handlers.InvokeRequest{
				Model:  "mock-test",
				Prompt: fmt.Sprintf("Generate code for task %d", index),
			}

			resp, err := invokeHandler.Handle(ctx, req)
			if err != nil {
				errors <- fmt.Errorf("invoke %d failed: %w", index, err)
				return
			}

			outputIDs <- resp.OutputID
		}(i)
	}

	// Wait for all to complete
	wg.Wait()
	close(errors)
	close(outputIDs)

	duration := time.Since(start)

	// Check for errors
	var errorCount int
	for err := range errors {
		t.Errorf("Concurrent invoke error: %v", err)
		errorCount++
	}

	// Collect output IDs
	var ids []string
	for id := range outputIDs {
		ids = append(ids, id)
	}

	// Verify results
	if errorCount > 0 {
		t.Errorf("Had %d errors out of %d concurrent calls", errorCount, numCalls)
	}

	if len(ids) != numCalls-errorCount {
		t.Errorf("Expected %d successful outputs, got %d", numCalls-errorCount, len(ids))
	}

	// Check that all IDs are unique
	uniqueIDs := make(map[string]bool)
	for _, id := range ids {
		if uniqueIDs[id] {
			t.Errorf("Duplicate output ID: %s", id)
		}
		uniqueIDs[id] = true
	}

	t.Logf("Completed %d concurrent invokes in %v (%.2f ops/sec)", 
		numCalls, duration, float64(numCalls)/duration.Seconds())
}

// TestConcurrentCheckCalls tests concurrent check operations
func TestConcurrentCheckCalls(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	store, err := storage.NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a test output first
	output := createTestOutput()
	if err := store.Save(output); err != nil {
		t.Fatalf("Failed to save test output: %v", err)
	}

	checkHandler := handlers.NewCheckHandler(store)
	ctx := context.Background()

	// Number of concurrent calls
	numCalls := 50
	var wg sync.WaitGroup
	errors := make(chan error, numCalls)

	// Track timing
	start := time.Now()

	// Launch concurrent checks
	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req := handlers.CheckRequest{
				OutputID: output.ID,
			}

			_, err := checkHandler.Handle(ctx, req)
			if err != nil {
				errors <- fmt.Errorf("check %d failed: %w", index, err)
			}
		}(i)
	}

	// Wait for all to complete
	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	var errorCount int
	for err := range errors {
		t.Errorf("Concurrent check error: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("Had %d errors out of %d concurrent calls", errorCount, numCalls)
	}

	t.Logf("Completed %d concurrent checks in %v (%.2f ops/sec)", 
		numCalls, duration, float64(numCalls)/duration.Seconds())
}

// TestConcurrentReadCalls tests concurrent read operations
func TestConcurrentReadCalls(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	store, err := storage.NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a test output with some content
	output := createTestOutput()
	output.Response.Raw = generateLargeContent(10000) // 10KB content
	if err := store.Save(output); err != nil {
		t.Fatalf("Failed to save test output: %v", err)
	}

	readHandler := handlers.NewReadHandler(store)
	ctx := context.Background()

	// Number of concurrent calls
	numCalls := 30
	var wg sync.WaitGroup
	errors := make(chan error, numCalls)

	// Track timing
	start := time.Now()

	// Launch concurrent reads with different extract options
	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Vary the extract option
			extractOptions := []string{"all", "code", "explanation"}
			extract := extractOptions[index%len(extractOptions)]

			req := handlers.ReadRequest{
				OutputID: output.ID,
				Options: handlers.ReadOptions{
					Extract: extract,
				},
			}

			_, err := readHandler.Handle(ctx, req)
			if err != nil {
				errors <- fmt.Errorf("read %d failed: %w", index, err)
			}
		}(i)
	}

	// Wait for all to complete
	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	var errorCount int
	for err := range errors {
		t.Errorf("Concurrent read error: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("Had %d errors out of %d concurrent calls", errorCount, numCalls)
	}

	t.Logf("Completed %d concurrent reads in %v (%.2f ops/sec)", 
		numCalls, duration, float64(numCalls)/duration.Seconds())
}

// TestMixedConcurrentOperations tests all three operations concurrently
func TestMixedConcurrentOperations(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	store, err := storage.NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	extractFactory := extractor.NewFactory()
	providerFactory := &mockProviderFactory{}
	
	invokeHandler := handlers.NewInvokeHandler(providerFactory, store, extractFactory)
	checkHandler := handlers.NewCheckHandler(store)
	readHandler := handlers.NewReadHandler(store)
	
	ctx := context.Background()

	// Create some initial outputs
	var initialOutputs []string
	for i := 0; i < 5; i++ {
		output := createTestOutput()
		output.ID = fmt.Sprintf("out_initial_%d", i)
		if err := store.Save(output); err != nil {
			t.Fatalf("Failed to save initial output: %v", err)
		}
		initialOutputs = append(initialOutputs, output.ID)
	}

	// Track operations
	type operation struct {
		opType string
		err    error
	}
	operations := make(chan operation, 100)

	var wg sync.WaitGroup
	start := time.Now()

	// Launch invoke operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			req := handlers.InvokeRequest{
				Model:  "mock-test",
				Prompt: fmt.Sprintf("Task %d", index),
			}
			
			_, err := invokeHandler.Handle(ctx, req)
			operations <- operation{"invoke", err}
		}(i)
	}

	// Launch check operations
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			// Use one of the initial outputs
			outputID := initialOutputs[index%len(initialOutputs)]
			req := handlers.CheckRequest{OutputID: outputID}
			
			_, err := checkHandler.Handle(ctx, req)
			operations <- operation{"check", err}
		}(i)
	}

	// Launch read operations
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			// Use one of the initial outputs
			outputID := initialOutputs[index%len(initialOutputs)]
			req := handlers.ReadRequest{
				OutputID: outputID,
				Options:  handlers.ReadOptions{Extract: "all"},
			}
			
			_, err := readHandler.Handle(ctx, req)
			operations <- operation{"read", err}
		}(i)
	}

	// Wait for all operations
	wg.Wait()
	close(operations)
	
	duration := time.Since(start)

	// Count results
	counts := make(map[string]int)
	errors := make(map[string]int)
	
	for op := range operations {
		counts[op.opType]++
		if op.err != nil {
			errors[op.opType]++
			t.Errorf("%s operation failed: %v", op.opType, op.err)
		}
	}

	// Report results
	t.Logf("Mixed concurrent operations completed in %v:", duration)
	for opType, count := range counts {
		errorCount := errors[opType]
		successRate := float64(count-errorCount) / float64(count) * 100
		t.Logf("  %s: %d operations, %d errors (%.1f%% success rate)",
			opType, count, errorCount, successRate)
	}
	
	totalOps := 0
	for _, count := range counts {
		totalOps += count
	}
	t.Logf("  Total: %.2f ops/sec", float64(totalOps)/duration.Seconds())
}

// Helper function to generate large content
func generateLargeContent(size int) string {
	content := "Here's a sample response with code:\n\n```python\ndef hello():\n    print('Hello, World!')\n```\n\n"
	
	// Repeat content to reach desired size
	for len(content) < size {
		content += "This is additional explanation text to fill up the content. "
	}
	
	return content[:size]
}

// BenchmarkInvokeHandler benchmarks the invoke handler
func BenchmarkInvokeHandler(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	store, _ := storage.NewFileStore(tempDir)
	extractFactory := extractor.NewFactory()
	providerFactory := &mockProviderFactory{}
	handler := handlers.NewInvokeHandler(providerFactory, store, extractFactory)
	
	ctx := context.Background()
	req := handlers.InvokeRequest{
		Model:  "mock-test",
		Prompt: "Generate a function",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Handle(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCheckHandler benchmarks the check handler
func BenchmarkCheckHandler(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	store, _ := storage.NewFileStore(tempDir)
	
	// Create test output
	output := createTestOutput()
	_ = store.Save(output)
	
	handler := handlers.NewCheckHandler(store)
	ctx := context.Background()
	req := handlers.CheckRequest{OutputID: output.ID}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Handle(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkReadHandler benchmarks the read handler
func BenchmarkReadHandler(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	store, _ := storage.NewFileStore(tempDir)
	
	// Create test output with content
	output := createTestOutput()
	output.Response.Raw = generateLargeContent(5000)
	_ = store.Save(output)
	
	handler := handlers.NewReadHandler(store)
	ctx := context.Background()
	req := handlers.ReadRequest{
		OutputID: output.ID,
		Options:  handlers.ReadOptions{Extract: "all"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Handle(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}