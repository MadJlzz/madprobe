// Service contains everything that relates to probe persistence.
// Validation is made on this layer too.
package service

import (
	"fmt"
	"sync"
)

// Once is an object that will perform exactly one action.
// It is used to ensure ProbeService is a singleton.
var once sync.Once

var instance *ProbeService

// Probe is the model required by the service to manipulate
// the resource.
type Probe struct {
	Name  string
	URL   string
	Delay uint
}

// ProbeService is an implementation
// of the interface controller.ProbeService
type ProbeService struct{}

// NewProbeService allow to create a new ProbeService
// implemented as a singleton.
func NewProbeService() *ProbeService {
	once.Do(func() {
		instance = &ProbeService{}
	})
	return instance
}

// Create does nothing but printing the given probe to the service layer.
// In the future, we should validate data received from the controller layer
// and find a way to persist it.
func (ps *ProbeService) Create(probe Probe) error {
	fmt.Printf("Probe: +%v has been successfully created.\n", probe)
	return nil
}
