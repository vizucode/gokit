package types

type PublisherArgument struct {
	PriorityMessage int
	CorrelationId   string
	Topic           string
	Exchange        string
	Queue           string
	Channel         string
	Key             string
	Headers         map[string]interface{}
	Message         []byte
}
