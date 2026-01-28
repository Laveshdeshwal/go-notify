package storage

type Database interface {
	Init() error
	GetDB() error
	Close() error
	CreateSchema() error
}
