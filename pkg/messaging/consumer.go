package messaging

type Consumer interface {
    Subscribe(map[string]func(event interface{}) error) error

    Close() error
}
