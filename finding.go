package pipeline

type Severity string

const (
	Low      Severity = "low"
	Medium   Severity = "medium"
	High     Severity = "high"
	Critical Severity = "critical"
)

type Topic string

const (
	TopicDefault  Topic = "findings.default"
	TopicPriority Topic = "findings.priority"
	TopicAlerts   Topic = "findings.alerts"
)
