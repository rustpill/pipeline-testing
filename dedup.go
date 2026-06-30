package pipeline

import "sync"

type Deduper struct {
	mu sync.Mutex
	// hashset
	seen map[string]struct{}
}

func NewDeduper() *Deduper {
	return &Deduper{seen: make(map[string]struct{})}
}

// keep mutations behind method
func (d *Deduper) Seen(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.seen[key]; ok {
		return true
	}
	d.seen[key] = struct{}{}
	return false
}
