package persistence

import (
	"errors"
	"io"
)

// Error returned when a connection wasn't established.
var ErrPersisterInitialization = errors.New("could not initialize connection")

// Error returned when an insertion fails.
var ErrPersisterInsertion = errors.New("could not insert entity")

// Error returned when a technical problem occurs during Get/GetAll
var ErrPersisterGet = errors.New("could not get entity(ies)")

// Error returned when a deletion fails.
var ErrPersisterDeletion = errors.New("could not delete entity")

// Any implementation that wishes to persist a probe
// must satisfy the following contract.
type Persister interface {
	Insert(entity *Entity) error
	Get(name string) (*Entity, error)
	GetAll() ([]*Entity, error)
	Delete(name string) error
}

// More specific type of an Persister that has to close one of it's resource.
type PersistCloser interface {
	Persister
	io.Closer
}

// Represent the data model that is stored in a file, database, etc...
type Entity struct {
	Name  string
	URL   string
	Delay uint
}

// Simple function that creates an entity given the parameters.
func NewEntity(name, URL string, delay uint) *Entity {
	return &Entity{
		Name:  name,
		URL:   URL,
		Delay: delay,
	}
}
