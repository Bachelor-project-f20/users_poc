package storage

type DbType int

const (
	MySQL DbType = iota
)

type Outbox interface {
	Insert(interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
}
