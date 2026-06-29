package pipeline

import "testing"

func TestRoute(t *testing.T) {

	tests := []struct {
		name string
		in   Severity
		want Topic
	}{
		{"critical routes to alerts", Critical, TopicAlerts},
		{"high routes to priority", High, TopicPriority},
		{"medium routes to default", Medium, TopicDefault},
		{"low routes to default", Low, TopicDefault},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Route(tt.in)
			if got != tt.want {
				t.Errorf("Route(%v) = %v, want %v", Critical, got, tt.want)
			}
		})
	}
}
