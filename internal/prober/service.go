// Service contains everything that relates to probe persistence.
// Validation is made on this layer too.
package prober

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	ErrProbeAlreadyExist = errors.New("probe with this name already exists")
	ErrProbeNotFound     = errors.New("probe was not found")
)

// Once is an object that will perform exactly one action.
// It is used to ensure ProbeService is a singleton.
var once sync.Once
var instance *service

// ProbeService is an implementation of the interface controller.ProbeService
type service struct {
	client   *http.Client
	probes   map[string]*Probe
	alertBus chan<- Probe
}

// NewProbeService allow to create a new probe service implemented as a singleton.
func NewProbeService(client *http.Client, alertBus chan<- Probe) *service {
	once.Do(func() {
		instance = &service{
			client:   client,
			probes:   make(map[string]*Probe),
			alertBus: alertBus,
		}
	})
	return instance
}

// Create does nothing but registering the given probe.
// Validation is made before storing the probe to be sure nothing partially configured enters the system.
func (ps *service) Create(probe Probe) error {
	err := runValidators(probe, nameInvalid, urlInvalid, delayInvalid)
	if err != nil {
		return err
	}
	if _, ok := ps.probes[probe.Name]; ok {
		return ErrProbeAlreadyExist
	}

	probe.finish = make(chan bool, 1)
	ps.probes[probe.Name] = &probe

	go ps.run(ps.probes[probe.Name])

	log.Printf("Probe [%s] has been successfuly created.\n", probe.Name)
	return nil
}

// Read retrieve a probe with the given name in the system.
// No validation is required. Returns the probe or ErrProbeNotFound is not probe has been found.
func (ps *service) Read(name string) (*Probe, error) {
	probe, ok := ps.probes[name]
	if !ok {
		return nil, ErrProbeNotFound
	}
	return probe, nil
}

// ReadAll retrieve all probes in the system.
func (ps *service) ReadAll() []*Probe {
	var probes []*Probe
	for _, value := range ps.probes {
		probes = append(probes, value)
	}
	return probes
}

// Update is a bit more complicated.
// Technically, it deletes the running probe and creates a new one with the given values.
func (ps *service) Update(name string, probe Probe) error {
	err := runValidators(probe, nameInvalid, urlInvalid, delayInvalid)
	if err != nil {
		return err
	}
	if _, ok := ps.probes[name]; !ok {
		return ErrProbeNotFound
	}
	if _, ok := ps.probes[probe.Name]; name != probe.Name && ok {
		return ErrProbeAlreadyExist
	}

	ps.probes[name].finish <- true
	delete(ps.probes, name)

	probe.finish = make(chan bool, 1)
	ps.probes[probe.Name] = &probe

	go ps.run(ps.probes[probe.Name])

	log.Printf("Probe [%s] has been successfuly updated.\n", probe.Name)
	return nil
}

// Delete erase an existing probe from the system.
// Validation is made before deletion to be sure nothing get removed by error.
func (ps *service) Delete(name string) error {
	probe := Probe{Name: name}
	err := runValidators(probe, nameInvalid)
	if err != nil {
		return err
	}

	mapProbe, ok := ps.probes[name]
	if !ok {
		return ErrProbeNotFound
	}

	mapProbe.finish <- true
	delete(ps.probes, mapProbe.Name)

	return nil
}

// run launch probes in a separate goroutine.
func (ps *service) run(probe *Probe) {
	var oldStatus string
	for {
		select {
		case <-probe.finish:
			log.Printf("<<HTTP PROBE [%s]>> Stopping probe...\n", probe.Name)
			return
		default:
			resp, err := ps.client.Get(probe.URL)
			oldStatus = probe.Status
			if err != nil {
				probe.Status = "DOWN"
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] is down.\n", probe.Name, probe.URL)
			} else if resp.StatusCode != 200 {
				b, _ := ioutil.ReadAll(resp.Body)
				probe.Status = "DOWN"
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] returned an error. got: ['%v']\n", probe.Name, probe.URL, string(b))
			} else {
				probe.Status = "UP"
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
