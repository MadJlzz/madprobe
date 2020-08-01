package persistence

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const probeBucket = "probe"

// Implementation of a Persister by using BoltDB
// as a key/value storage.
type boltDBClient struct {
	boltDB *bolt.DB
}

// Initialize a new BoltDB connection, creates the main bucket, etc...
// Returns a nil connection if there was an error.
// Also returns an error that should be considered.
func NewBoltDBClient(filepath string) (*boltDBClient, error) {
	con, err := bolt.Open(filepath, 0600, bolt.DefaultOptions)
	if err != nil {
		return nil, errors.Wrap(err, ErrPersisterInitialization.Error())
	}
	err = con.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(probeBucket))
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, ErrPersisterInitialization.Error())
	}
	return &boltDBClient{boltDB: con}, nil
}

func (c *boltDBClient) Close() error {
	return c.boltDB.Close()
}

// Insert a new entity inside BoltDB. Returns nil if there was no errors.
func (c *boltDBClient) Insert(entity *Entity) error {
	err := c.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		bytes, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(entity.Name), bytes)
	})
	return errors.Wrap(err, ErrPersisterInsertion.Error())
}

// Delete delete probe by Name, returns nil error on success.
func (c *boltDBClient) Delete(name string) error {
	err := c.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		return bucket.Delete([]byte(name))
	})
	return errors.Wrap(err, ErrPersisterDeletion.Error())
}

// Get returns an entity with it's name.
// Return value can be nil for the entity is nothing is found.
// An error can be returned if a technical issue occurred.
func (c *boltDBClient) Get(name string) (*Entity, error) {
	var entity Entity
	err := c.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		probeBytes := bucket.Get([]byte(name))
		if len(probeBytes) == 0 {
			return nil
		}
		return json.Unmarshal(probeBytes, &entity)
	})
	return &entity, errors.Wrap(err, ErrPersisterGet.Error())
}

// GetAll returns all entities from the database or an empty slice if nothing actually stored.
// An error is returned if any technical error occurs.
func (c *boltDBClient) GetAll() ([]*Entity, error) {
	var err error
	var entities []*Entity
	err = c.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		return bucket.ForEach(func(_, data []byte) error {
			var entity Entity
			jsonErr := json.Unmarshal(data, &entity)
			if jsonErr != nil {
				return errors.Wrap(err, jsonErr.Error())
			}
			entities = append(entities, &entity)
			return nil
		})
	})
	return entities, err
}
