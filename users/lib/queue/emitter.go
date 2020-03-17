package queue

import ob "github.com/dueruen/go-outbox"

type EventEmitter interface {
	Emit(e ob.Event) error
}
