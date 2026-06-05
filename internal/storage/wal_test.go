package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestData represents a simple struct we want to store and retrieve.
type TestData struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func TestWalStorage_BasicOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wal-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	snapshotPath := filepath.Join(tempDir, "test-snapshot.json")
	walPath := filepath.Join(tempDir, "test.wal")

	// Initialize WalStorage with a short compaction interval
	store, err := NewWalStorage[string, TestData](snapshotPath, walPath, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to initialize WalStorage: %v", err)
	}

	ctx := context.Background()

	// 1. Put data
	data1 := TestData{ID: "1", Value: "hello"}
	err = store.Put(ctx, "key1", data1)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// 2. Get data
	retrieved, exists, err := store.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !exists {
		t.Fatalf("Expected key1 to exist")
	}
	if retrieved.Value != "hello" {
		t.Errorf("Expected 'hello', got '%s'", retrieved.Value)
	}

	// 3. Delete data
	err = store.Delete(ctx, "key1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, exists, err = store.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected key1 to be deleted")
	}
}

func TestWalStorage_PersistenceAndCompaction(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wal-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	snapshotPath := filepath.Join(tempDir, "test-snapshot.json")
	walPath := filepath.Join(tempDir, "test.wal")

	// Phase 1: Write data and let it compact
	store1, _ := NewWalStorage[string, TestData](snapshotPath, walPath, 50*time.Millisecond)
	ctx := context.Background()
	
	_ = store1.Put(ctx, "key1", TestData{ID: "1", Value: "persisted data"})
	
	// Wait for compaction (snapshot creation)
	time.Sleep(200 * time.Millisecond)

	// Phase 2: Create a NEW storage instance pointing to the same files to test recovery
	store2, err := NewWalStorage[string, TestData](snapshotPath, walPath, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to load existing WalStorage: %v", err)
	}

	retrieved, exists, err := store2.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !exists {
		t.Fatalf("Data was not persisted across instances")
	}
	if retrieved.Value != "persisted data" {
		t.Errorf("Expected 'persisted data', got '%s'", retrieved.Value)
	}
}
