package persistence

import (
	"encoding/json"
	"fmt"
	"github.com/madjlzz/madprobe/internal/prober"

	"github.com/boltdb/bolt"
)

type probeEntity struct {
	Name  string
	URL   string
	Delay uint
}

func (p *probeEntity) toProbe() *prober.Probe {
	return &prober.Probe{
		Name:  p.Name,
		URL:   p.URL,
		Delay: p.Delay,
	}
}

func newProbeEntity(probe *prober.Probe) *probeEntity {
	// Status and Finish chan not persisted
	return &probeEntity{
		Name:  probe.Name,
		URL:   probe.URL,
		Delay: probe.Delay,
	}
}

// InsertProbe returns nil error on success.
func (c *PersistenceClient) InsertProbe(probe *prober.Probe) error {
	var probeEntity = newProbeEntity(probe)
	err := c.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		bytes, err := json.Marshal(probeEntity)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(probeEntity.Name), bytes)
	})
	if err != nil {
		return fmt.Errorf("failed to insert probe: %w", err)
	}
	probe.Finish = make(chan bool, 1)
	c.probeFinishChannels[probe.Name] = probe.Finish
	return nil
}

// DeleteProbe delete probe by Name, returns nil error on success.
func (c *PersistenceClient) DeleteProbeByName(probeName string) error {
	err := c.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		return bucket.Delete([]byte(probeName))
	})
	if err != nil {
		return fmt.Errorf("failed to delete probe: %w", err)
	}
	if finishChannel, exist := c.probeFinishChannels[probeName]; exist {
		finishChannel <- true
		delete(c.probeFinishChannels, probeName)
	}
	return nil
}

// ExistByName returns true if probe with this name exist
func (c *PersistenceClient) ExistProbeByName(probeName string) (bool, error) {
	probe, err := c.GetProbe(probeName)
	return probe != nil, err
}

// GetProbe returns a probe by name, result can be nil if no probe is found
func (c *PersistenceClient) GetProbe(probeName string) (*prober.Probe, error) {
	var probeEntity probeEntity
	err := c.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		probeBytes := bucket.Get([]byte(probeName))
		if probeBytes == nil {
			return nil
		}
		err := json.Unmarshal(probeBytes, &probeEntity)
		if err != nil {
			return fmt.Errorf("failed to unmarshal probe: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve probe: %w", err)
	}
	if probeEntity.Name == "" {
		return nil, nil // No probe found
	}
	return probeEntity.toProbe(), nil
}

// GetAllProbe returns all probes in database
func (c *PersistenceClient) GetAllProbes() ([]*prober.Probe, error) {
	var probes []*prober.Probe
	err := c.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(probeBucket))
		return bucket.ForEach(func(_, probeStored []byte) error {
			var pEntity probeEntity
			err := json.Unmarshal(probeStored, &pEntity)
			if err != nil {
				return err
			}
			probes = append(probes, pEntity.toProbe())
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read all probes: %w", err)
	}
	return probes, nil
}
