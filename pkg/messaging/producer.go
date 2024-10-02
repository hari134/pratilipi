package messaging

type Producer interface {
    Emit(topic string, event interface{}) error

    Close() error
}
