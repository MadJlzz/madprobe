package persistence

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/madjlzz/madprobe/internal/model"
)

type probeEntity struct {
	Name  string
	URL   string
	Delay uint
}

func newProbeEntity(probe *model.Probe) *probeEntity {
	// Status and Finish chan not persisted
	return &probeEntity{
		Name:  probe.Name,
		URL:   probe.URL,
		Delay: probe.Delay,
	}
}

// InsertProbe returns nil error on success.
func (c *PersistenceClient) InsertProbe(probe *model.Probe) error {
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
	c.probes[probe.Name] = probe
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
	if probe, exist := c.probes[probeName]; exist {
		probe.Finish <- true
		delete(c.probes, probe.Name)
	}
	return nil
}

// ExistByName returns true if probe with this name exist
func (c *PersistenceClient) ExistProbeByName(probeName string) bool {
	_, exist := c.probes[probeName]
	return exist
}

// GetProbe returns a probe by name
func (c *PersistenceClient) GetProbe(probeName string) (*model.Probe, bool) {
	probe, ok := c.probes[probeName]
	return probe, ok
}

// GetAllProbe returns all probes
func (c *PersistenceClient) GetAllProbe() []*model.Probe {
	var probes []*model.Probe
	for _, probe := range c.probes {
		probes = append(probes, probe)
	}
	return probes
}

// readAllProbes returns all probes in database
func (c *PersistenceClient) readAllStoredProbes() ([]*model.Probe, error) {
	var probes []*model.Probe
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

func (p *probeEntity) toProbe() *model.Probe {
	return &model.Probe{
		Name:  p.Name,
		URL:   p.URL,
		Delay: p.Delay,
	}
}
