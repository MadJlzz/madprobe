package prober

// ProbeService represent the interface used to manipulate probes.
type ProbeService interface {
	Create(probe Probe) error
	Read(name string) (*Probe, error)
	ReadAll() ([]*Probe, error)
	Update(name string, probe Probe) error
	Delete(name string) error
}

// Probe is the model required by the service to manipulate the resource.
type Probe struct {
	Name   string
	URL    string
	Status string
	Delay  uint
	Finish ProbeFinishChannel
}

type ProbeFinishChannel = chan bool
