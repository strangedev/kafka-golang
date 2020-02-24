package consumer

// Consumer is the generic interface for all consumers.
type Consumer interface {
	// Run starts the consumer. It is necessary to call this method, otherwise the Consumer
	// won't consume any events.
	Run() (stop chan bool, err error)
}
