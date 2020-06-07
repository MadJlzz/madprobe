package model

// Probe is the model required by the service to manipulate
// the resource.
type Probe struct {
	Name   string
	URL    string
	Status string
	Delay  uint
	Finish chan bool
}
