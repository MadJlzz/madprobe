package prober

// ProbeService represent the interface used to manipulate probes.
type ProbeService interface {
	Insert(probe Probe) error
	Get(name string) (*Probe, error)
	GetAll() ([]*Probe, error)
	Delete(name string) error
}

// TODO: we need a solution to decouple Run() from the package prober so that it becomes independent.
// ProbeRunner contract purpose is to run a probe.
type ProbeRunner interface {
	Run(probe *Probe)
}

// Probe is the model required by the service to manipulate the resource.
type Probe struct {
	Name   string
	URL    string
	Status string
	Delay  uint
	Finish chan bool
}

// Creates a new Probe with the given parameters.
func NewProbe(name, URL string, delay uint) *Probe {
	return &Probe{
		Name:   name,
		URL:    URL,
		Delay:  delay,
		Finish: make(chan bool),
	}
}
