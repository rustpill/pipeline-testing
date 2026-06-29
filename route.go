package pipeline

func Route(s Severity) (topic Topic) {
	switch s {
	case Critical:
		topic = TopicAlerts
	case High:
		topic = TopicPriority
	default:
		topic = TopicDefault
	}
	return
}
