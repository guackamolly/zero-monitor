package db

// Simple abstraction of a Database "object"
type Database interface {
	Open() error
	Tables() []Table
	Table(id string) (Table, bool)
	Close() error
}
