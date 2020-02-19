package consumer

type Consumer interface {
	Run() (stop chan bool, err error)
}
