package pipeline

import "fmt"

// Pipeline
type Pipeline struct {
	out        Producer
	deadLetter Producer
	dedup      *Deduper
}

func NewPipeline(out, deadLetter Producer) *Pipeline {
	return &Pipeline{
		out:        out,
		deadLetter: deadLetter,
		dedup:      NewDeduper(),
	}
}

func (p *Pipeline) Handle(raw []byte) error {
	f, err := Decode(raw)

	if err != nil {
		return p.toDeadLetter(raw, err)
	}

	if err := f.Validate(); err != nil {
		return p.toDeadLetter(raw, err)
	}

	// check if seen
	if p.dedup.Seen(f.ID) {
		return nil
	}

	// Create Message from bytes
	msg := Message{
		// Use Route to establish which Topic to "send to"
		Topic: Route(f.Severity),
		Key:   f.ID,
		Value: raw,
	}

	return p.out.Produce(msg)
}

func (p *Pipeline) toDeadLetter(raw []byte, cause error) error {
	msg := Message{Topic: TopicDeadLetter, Value: raw}

	if err := p.deadLetter.Produce(msg); err != nil {
		return fmt.Errorf("dead-letter produce failed (cause: %v): %w", cause, err)
	}

	return nil
}
