package pipeline

import "testing"

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
		if got := out.Produced[0].Topic; got != TopicAlerts {
			t.Errorf("produced to %v, want %v", got, TopicAlerts)
		}
		if got := out.Produced[0].Key; got != "CVE-2021-44228" {
			t.Errorf("message key = %q, want the finding id", got)
		}
		if len(dlq.Produced) != 0 {
			t.Errorf("expected nothing on the dead-letter topic, got %d", len(dlq.Produced))
		}

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
			if len(out.Produced) != 0 {
				t.Errorf("bad input should not reach the out topic, got %d messages", len(out.Produced))
			}
			if len(dlq.Produced) != 1 {
				t.Fatalf("expected 1 message on the dead-letter topic, got %d", len(dlq.Produced))
			}
			if got := dlq.Produced[0].Topic; got != TopicDeadLetter {
				t.Errorf("dead-lettered to %v, want %v", got, TopicDeadLetter)
			}

		})
	}
}
