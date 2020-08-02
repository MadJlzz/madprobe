// Service contains everything that relates to probe persistence.
// Validation is made on this layer too.
package prober

import (
	"errors"
	"github.com/madjlzz/madprobe/internal/persistence"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const downStatus = "DOWN"
const upStatus = "UP"

var (
	ErrProbeAlreadyExist = errors.New("probe with this name already exists")
	ErrProbeNotFound     = errors.New("probe was not found")
)

var instance *service

// ProbeService is an implementation of the interface controller.ProbeService
type service struct {
	client    *http.Client
	persister persistence.Persister
	alertBus  chan<- Probe
	probes    map[string]*Probe
}

// NewProbeService allow to create a new probe service implemented as a singleton.
func NewProbeService(httpClient *http.Client, persister persistence.Persister, alertBus chan<- Probe) *service {
	instance = &service{
		client:    httpClient,
		persister: persister,
		alertBus:  alertBus,
		probes:    make(map[string]*Probe),
	}

	err := instance.runProbes()
	if err != nil {
		log.Fatalf("failed to start probe service: %s", err.Error())
	}
	return instance
}

// Insert does nothing but registering the given probe.
// Validation is made before storing the probe to be sure nothing partially configured enters the system.
// Local cache is also updated.
func (ps *service) Insert(probe Probe) error {
	err := runValidators(probe, nameInvalid, urlInvalid, delayInvalid)
	if err != nil {
		return err
	}

	entity, err := ps.persister.Get(probe.Name)
	if err != nil {
		return err
	}
	if entity.Name != "" {
		return ErrProbeAlreadyExist
	}

	entity = persistence.NewEntity(probe.Name, probe.URL, probe.Delay)
	err = ps.persister.Insert(entity)
	if err != nil {
		return err
	}

	ps.probes[probe.Name] = &probe
	go ps.run(&probe)

	log.Printf("Probe [%s] has been successfuly created.\n", probe.Name)
	return nil
}

// Get retrieve a probe with the given name in the system.
// No validation is required. Returns the probe or ErrProbeNotFound if no probe has been found.
func (ps *service) Get(name string) (*Probe, error) {
	probe, err := ps.persister.Get(name)
	if err != nil {
		return nil, err
	}
	if probe == nil {
		return nil, ErrProbeNotFound
	}
	return ps.probes[name], nil
}

// GetAll retrieve all probes in the system or an empty slice.
func (ps *service) GetAll() ([]*Probe, error) {
	entities, err := ps.persister.GetAll()
	if err != nil {
		return nil, err
	}
	var probes []*Probe
	for _, entity := range entities {
		probes = append(probes, ps.probes[entity.Name])
	}
	return probes, nil
}

// Delete erase an existing probe from the system.
// Validation is made before deletion to be sure nothing get removed by error.
// Local cache is also updated.
func (ps *service) Delete(name string) error {
	probe := Probe{Name: name}
	err := runValidators(probe, nameInvalid)
	if err != nil {
		return err
	}

	err = ps.persister.Delete(name)
	if err != nil {
		return err
	}
	ps.probes[name].Finish <- true
	delete(ps.probes, name)

	return nil
}

// run launches probes in a separate goroutine.
func (ps *service) run(probe *Probe) {
	var oldStatus string
	for {
		select {
		case <-probe.Finish:
			log.Printf("<<HTTP PROBE [%s]>> Stopping probe...\n", probe.Name)
			return
		default:
			resp, err := ps.client.Get(probe.URL)
			oldStatus = probe.Status
			if err != nil {
				probe.Status = downStatus
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] is down.\n", probe.Name, probe.URL)
			} else if resp.StatusCode != 200 {
				b, _ := ioutil.ReadAll(resp.Body)
				probe.Status = downStatus
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] returned an error. got: ['%v']\n", probe.Name, probe.URL, string(b))
			} else {
				probe.Status = upStatus
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] is alive.\n", probe.Name, probe.URL)
				_ = resp.Body.Close()
			}
		}
		// If the status has changed, we can send an event to the alerter bus...
		if oldStatus != probe.Status {
			ps.alertBus <- *probe
		}
		time.Sleep(time.Duration(probe.Delay) * time.Second)
	}
}

func (ps *service) runProbes() error {
	entities, err := ps.persister.GetAll()
	if err != nil {
		return err
	}
	for _, entity := range entities {
		probe := NewProbe(entity.Name, entity.URL, entity.Delay)
		ps.probes[entity.Name] = probe
		go ps.run(probe)
	}
	return nil
}
