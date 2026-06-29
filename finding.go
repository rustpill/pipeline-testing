package pipeline

import "errors"

type Severity string

var (
	ErrEmptyID         = errors.New("finding: empty id")
	ErrUnknownSeverity = errors.New("finding: unknown severity")
	ErrCVSSOutOfRange  = errors.New("finding: cvss out of range")
)

const (
	Low      Severity = "low"
	Medium   Severity = "medium"
	High     Severity = "high"
	Critical Severity = "critical"
)

func (s Severity) valid() bool {
	switch s {
	case Low, Medium, High, Critical:
		return true
	default:
		return false
	}
}

type Topic string

const (
	TopicDefault  Topic = "findings.default"
	TopicPriority Topic = "findings.priority"
	TopicAlerts   Topic = "findings.alerts"
)

type Finding struct {
	ID       string   `json:"id"`
	Package  string   `json:"package"`
	Severity Severity `json:"severity"`
	CVSS     float64  `json:"cvss"`
}

func (f Finding) Validate() error {
	if f.ID == "" {
		return ErrEmptyID
	}

	if !f.Severity.valid() {
		return ErrUnknownSeverity
	}

	if f.CVSS < 0 || f.CVSS > 10 {
		return ErrCVSSOutOfRange
	}
	return nil
}
