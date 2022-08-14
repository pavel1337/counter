package mapper

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	oneSec = time.Second
)

func TestFileStore(t *testing.T) {
	path, err := newTempFileName()
	if err != nil {
		t.Fatal(err)
	}

	store, err := New(path)
	if err != nil {
		t.Fatalf("Error creating store: %s", err)
	}

	for i := 0; i < 1000; i++ {
		if err := store.Add(time.Now().Add(-oneSec)); err != nil {
			t.Fatalf("Error adding time: %s", err)
		}
	}

	last, err := store.Last(2)
	if err != nil {
		t.Fatalf("Error getting last time: %s", err)
	}
	if last != 1000 {
		t.Fatalf("expected 1000, got %d", last)
	}

	store, err = New(path)
	if err != nil {
		t.Fatalf("Error creating store: %s", err)
	}

	last, err = store.Last(2)
	if err != nil {
		t.Fatalf("Error getting last time: %s", err)
	}
	if last != 1000 {
		t.Fatalf("expected 1000, got %d", last)
	}

	last, err = store.Last(1)
	if err != nil {
		t.Fatalf("Error getting last time: %s", err)
	}

	if last != 0 {
		t.Fatalf("expected 0, got %d", last)
	}

}

// BenchmarkStore benchmarks Store.
func BenchmarkStore(b *testing.B) {
	path, err := newTempFileName()
	if err != nil {
		b.Fatalf("Error creating temp file: %s", err)
	}

	store, err := New(path)
	if err != nil {
		b.Fatalf("Error creating store: %s", err)
	}

	for i := 0; i < b.N; i++ {
		if err := store.Add(time.Now()); err != nil {
			b.Fatalf("Error adding time: %s", err)
		}
	}
}

// newTempFileName returns a new temporary file name.
func newTempFileName() (string, error) {
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %s", err)
	}
	defer tempFile.Close()

	return tempFile.Name(), nil
}
