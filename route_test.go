package pipeline

import (
	"errors"
	"testing"
)

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

func TestFindingValidate(t *testing.T) {
	tests := []struct {
		name    string
		finding Finding
		wantErr error
	}{
		{
			name:    "valid finding",
			finding: Finding{ID: "CVE-2021-44228", Severity: Critical, CVSS: 10.0},
			wantErr: nil,
		},
		{
			name:    "empty id",
			finding: Finding{Severity: High, CVSS: 7.5},
			wantErr: ErrEmptyID,
		},
		{
			name:    "unknown severity",
			finding: Finding{ID: "x", Severity: "banana", CVSS: 5.0},
			wantErr: ErrUnknownSeverity,
		},
		{
			name:    "cvss too high",
			finding: Finding{ID: "x", Severity: Low, CVSS: 11.0},
			wantErr: ErrCVSSOutOfRange,
		},
		{
			name:    "cvss negative",
			finding: Finding{ID: "x", Severity: Low, CVSS: -1},
			wantErr: ErrCVSSOutOfRange,
		},
		{
			name:    "empty id and bad cvss",
			finding: Finding{Severity: High, CVSS: 99},
			wantErr: ErrCVSSOutOfRange, // errors.Is still finds it inside the join
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.finding.Validate()

			if tt.wantErr == nil {
				// should return nil
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			// check if error matches expected
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want one matching %v", err, tt.wantErr)
			}

		})
	}
}
