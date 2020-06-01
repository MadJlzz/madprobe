package model

// Probe is the model required by the service to manipulate
// the resource.
type Probe struct {
	Name   string
	URL    string
	Status string
	Delay  uint
	Finish chan bool
	ID     int
}

func (p Probe) Copy() Probe {
	return Probe{
		Name:   p.Name,
		URL:    p.URL,
		Status: p.Status,
		Delay:  p.Delay,
		Finish: p.Finish,
		ID:     p.ID,
	}
}
