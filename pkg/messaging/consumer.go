package messaging

type Consumer interface {
    Subscribe(topic string, eventStruct interface{}, handler func(event interface{}) error) error

    Close() error
}
