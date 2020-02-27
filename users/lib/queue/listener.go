package queue

type EventListener interface {
	Listen(events ...string) (<-chan Event, <-chan error, error)
}
