// Service contains everything that relates to probe persistence.
// Validation is made on this layer too.
package service

import (
	"errors"
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

var instance *ProbeService

// Probe is the model required by the service to manipulate
// the resource.
type Probe struct {
	Name   string
	URL    string
	Delay  uint
	finish chan bool
}

// ProbeService is an implementation
// of the interface controller.ProbeService
type ProbeService struct {
	probes map[string]Probe
}

// NewProbeService allow to create a new ProbeService
// implemented as a singleton.
func NewProbeService() *ProbeService {
	once.Do(func() {
		instance = &ProbeService{
			probes: make(map[string]Probe),
		}
	})
	return instance
}

// Create does nothing but registering the given probe.
// Validation is made before storing the probe to be sure nothing partially configured enters the system.
func (ps *ProbeService) Create(probe Probe) error {
	err := runValidators(probe, nameInvalid, urlInvalid, delayInvalid)
	if err != nil {
		return err
	}
	if _, ok := ps.probes[probe.Name]; ok {
		return ErrProbeAlreadyExist
	}

	probe.finish = make(chan bool, 1)
	ps.probes[probe.Name] = probe

	go run(probe)

	log.Printf("Probe [%s] has been successfuly created.\n", probe.Name)
	return nil
}

// Read retrieve a probe with the given name in the system.
// No validation is required. Returns the probe or ErrProbeNotFound is not probe has been found.
func (ps *ProbeService) Read(name string) (*Probe, error) {
	probe, ok := ps.probes[name]
	if !ok {
		return nil, ErrProbeNotFound
	}
	return &probe, nil
}

// Delete erase an existing probe from the system.
// Validation is made before deletion to be sure nothing get removed by error.
func (ps *ProbeService) Delete(name string) error {
	probe := Probe{Name: name}
	err := runValidators(probe, nameInvalid)
	if err != nil {
		return err
	}

	probe, ok := ps.probes[name]
	if !ok {
		return ErrProbeNotFound
	}

	probe.finish <- true
	delete(ps.probes, probe.Name)

	return nil
}

func run(probe Probe) {
	for {
		select {
		case <-probe.finish:
			log.Printf("<<HTTP PROBE [%s]>> Stopping probe...\n", probe.Name)
			return
		default:
			resp, err := http.Get(probe.URL)
			if err != nil {
				log.Printf("<<HTTP PROBE [%s]>> Service targeting [%s] is down.\n", probe.Name, probe.URL)
			} else {
				log.Printf("<<HTTP PROBE [%s]>> Service targeting [%s] is alive.\n", probe.Name, probe.URL)
				_ = resp.Body.Close()
			}
		}
		time.Sleep(time.Duration(probe.Delay) * time.Second)
	}
}
