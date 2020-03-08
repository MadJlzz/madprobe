// Service contains everything that relates to probe persistence.
// Validation is made on this layer too.
package service

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Once is an object that will perform exactly one action.
// It is used to ensure ProbeService is a singleton.
var once sync.Once

var instance *ProbeService

// Probe is the model required by the service to manipulate
// the resource.
type Probe struct {
	Name     string
	URL      string
	Delay    uint
	receiver chan int
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

// Create does nothing but printing the given probe to the service layer.
// In the future, we should validate data received from the controller layer
// and find a way to persist it.
func (ps *ProbeService) Create(probe Probe) error {
	err := runValidators(probe, nameInvalid, urlInvalid, delayInvalid)
	if err != nil {
		return err
	}

	probe.receiver = make(chan int, 1)
	ps.probes[probe.Name] = probe

	go run(probe)

	log.Printf("Probe [%s] has been successfuly created.\n", probe.Name)
	return nil
}

func run(probe Probe) {
	for {
		select {
		case <-probe.receiver:
			log.Printf("<<HTTP PROBE [%s]>> stopping probe...\n", probe.Name)
			os.Exit(0)
		default:
			resp, err := http.Get(probe.URL)
			if err != nil {
				log.Printf("<<HTTP PROBE [%s]>> an error occured while doing HTTP call to [%s]. got '%s'\n", probe.Name, probe.URL, err)
			} else {
				log.Printf("<<HTTP PROBE [%s]>> Service is alive.\n", probe.Name)
				_ = resp.Body.Close()
			}
		}
		time.Sleep(time.Duration(probe.Delay) * time.Second)
	}
}
