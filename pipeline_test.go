package pipeline

import (
	"errors"
	"testing"
)

func TestPipelineHandle(t *testing.T) {
	t.Run("valid critical finding is produced to the alerts topic", func(t *testing.T) {
		out := &FakeProducer{}
		dlq := &FakeProducer{}
		p := NewPipeline(out, dlq)

		raw := []byte(`{"id":"CVE-2021-44228","severity":"critical","cvss":10}`)

		if err := p.Handle(raw); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out.Produced) != 1 {
			t.Fatalf("expected 1 message on out, got %d", len(out.Produced))
		}
		assertEqual(t, out.Produced[0].Topic, TopicAlerts, "Message should be sent to TopicAlerts topic:")

		assertEqual(t, out.Produced[0].Key, "CVE-2021-44228", "Message key should be finding ID:")

		assertEqual(t, len(dlq.Produced), 0, "DLQ should be empty")

	})
}

func TestPipelineDeadLetters(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{"malformed json", `{not json`},
		{"invalid finding (empty id)", `{"severity":"high","cvss":5}`},
		{"invalid finding (bad cvss)", `{"id":"x","severity":"low","cvss":42}`},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			out := &FakeProducer{}
			dlq := &FakeProducer{}
			p := NewPipeline(out, dlq)

			if err := p.Handle([]byte(test.raw)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertEqual(t, len(out.Produced), 0, "Bad input should not reach out topic:")

			if len(dlq.Produced) != 1 {
				t.Fatalf("expected 1 message on the dead-letter topic, got %d", len(dlq.Produced))
			}

			assertEqual(t, dlq.Produced[0].Topic, TopicDeadLetter, "TopicDeadLetter expected in DLQ:")

		})
	}
}

func TestDeadLetterFailsError(t *testing.T) {
	out := &FakeProducer{}
	dlq := &FakeProducer{
		// force error
		ProducedFn: func(Message) error {
			return errors.New("dead-letter broker unreachable")
		},
	}
	p := NewPipeline(out, dlq)

	err := p.Handle([]byte(`{not json`))
	if err == nil {
		t.Fatal("expected an error when the dead-letter producer fails, got nil")
	}
}

func TestPipelineDropsDuplicates(t *testing.T) {
	out := &FakeProducer{}
	dlq := &FakeProducer{}
	p := NewPipeline(out, dlq)

	raw := []byte(`{"id":"CVE-2021-44228","severity":"critical","cvss":10}`)

	_ = p.Handle(raw)
	_ = p.Handle(raw)

	assertEqual(t, len(out.Produced), 1, "Duplicates should be ignore:")

}
