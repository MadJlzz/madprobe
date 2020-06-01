package database

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	databaseFile = "madprobe.db"
	probeBucket  = "Probe"
)

type client struct {
	db *bolt.DB
}

var c *client

func GetClient() (*client, error) {
	if c != nil {
		return c, nil
	}
	var err error
	c, err = newClient()
	return c, err
}

func newClient() (*client, error) {
	db, err := bolt.Open(databaseFile, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bold database (%s): %w", databaseFile, err)
	}
	c := &client{
		db: db,
	}
	err = c.initBuckets()
	return c, err
}

func (c *client) initBuckets() error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(probeBucket))
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to create buckets if needed: %w", err)
	}
	return nil
}
