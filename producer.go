package pipeline

// Custom producer Message
type Message struct {
	Topic Topic
	Key   string
	Value []byte
}

type Producer interface {
	Produce(msg Message) error
}
