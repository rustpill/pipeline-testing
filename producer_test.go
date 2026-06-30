package pipeline

import "sync"

// Fake produce with custom behaviour for testing
type FakeProducer struct {
	mu sync.Mutex
	// Our "pipeline"
	Produced   []Message
	ProducedFn func(Message) error
}

func (f *FakeProducer) Produce(msg Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// handle func
	if f.ProducedFn != nil {
		if err := f.ProducedFn(msg); err != nil {
			return err // failed publish
		}
	}

	f.Produced = append(f.Produced, msg)
	return nil
}
