package pipeline

// Pipeline
type Pipeline struct {
	out        Producer
	deadLetter Producer
}

func NewPipeline(out, deadLetter Producer) *Pipeline {
	return &Pipeline{out: out, deadLetter: deadLetter}
}

func (p *Pipeline) Handle(raw []byte) error {
	f, err := Decode(raw)

	if err != nil {
		return err
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
