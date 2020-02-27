package queue

type Event interface {
	GetID() string
	GetPublisher() string
	GetTimestamp() int64
}
