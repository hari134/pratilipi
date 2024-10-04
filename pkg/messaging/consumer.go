package messaging

type Consumer interface {
    Subscribe(topic string,eventName string, handler func(event interface{}) error) error

    Close() error
}
