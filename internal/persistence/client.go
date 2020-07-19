package persistence

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/madjlzz/madprobe/internal/prober"
)

const (
	probeBucket = "probe"
)

type PersistenceClient struct {
	probeFinishChannels map[string]prober.ProbeFinishChannel
	boltDB              *bolt.DB
}

func NewPersistenceClient() (*PersistenceClient, error) {
	boltDB, err := newBoltDBClient()
	if err != nil {
		return nil, err
	}
	persistenceClient := &PersistenceClient{
		probeFinishChannels: make(map[string]prober.ProbeFinishChannel),
		boltDB:              boltDB,
	}
	return persistenceClient, nil
}

func (c *PersistenceClient) LoadProbes() ([]*prober.Probe, error) {
	probes, err := c.GetAllProbes()
	if err != nil {
		return nil, err
	}
	for _, probe := range probes {
		probe.Finish = make(chan bool, 1)
		c.probeFinishChannels[probe.Name] = probe.Finish
	}
	return probes, nil
}

func newBoltDBClient() (*bolt.DB, error) {
	// TODO: be able to change this thing.
	databaseFile := "madprobe.db"
	db, err := bolt.Open(databaseFile, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bold database (%s): %w", databaseFile, err)
	}
	err = initBuckets(db)
	return db, err
}

func initBuckets(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(probeBucket))
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to create buckets if needed: %w", err)
	}
	return nil
}
