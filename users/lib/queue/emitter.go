package queue

type EventEmitter interface {
	Emit(e Event) error
}
