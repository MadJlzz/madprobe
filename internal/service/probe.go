package service

import "fmt"

type Probe struct {
	Name  string
	URL   string
	Delay uint
}

// Create does nothing but printing the given probe to the service layer.
// In the future, we should validate data received from the controller layer
// and find a way to persist it.
func Create(probe Probe) error {
	fmt.Printf("Probe: +%v has been successfully created.\n", probe)
	return nil
}
