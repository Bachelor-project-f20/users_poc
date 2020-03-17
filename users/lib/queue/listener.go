package queue

import ob "github.com/dueruen/go-outbox"

type EventListener interface {
	Listen(events ...string) (<-chan ob.Event, <-chan error, error)
}
