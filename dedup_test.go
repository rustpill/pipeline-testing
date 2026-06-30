package pipeline

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestDeduper(t *testing.T) {
	d := NewDeduper()

	// .Seen() returns false if seen
	if d.Seen("finding-1") {
		t.Error("first sighting of finding-1 should report not-seen")
	}
	if !d.Seen("finding-1") {
		t.Error("second sighting of finding-1 should report seen")
	}
	if d.Seen("finding-2") {
		t.Error("first sighting of finding-2 should report not-seen")
	}
}

func TestDedupConcurrency(t *testing.T) {
	d := NewDeduper()

	const goroutines = 100
	var firstTimers atomic.Int64
	var wg sync.WaitGroup

	for range goroutines {
		wg.Go(func() {
			if !d.Seen("same-key") {
				firstTimers.Add(1)
			}
		})
	}
	wg.Wait()

	assertEqual(t, firstTimers.Load(), 1, "Should have exactly 1 seen:")
}
