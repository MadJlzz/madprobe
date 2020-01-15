package probe

import (
	"log"
)

type Prober interface {
	Probe()
}

type Probes struct {
	Pid  []Pid  `yaml:"pid"`
	Http []Http `yaml:"http"`

	probers []Prober
}

func Run(p *Probes) {
	p.registerProbes()
	log.Printf("Running probes defined in yaml configuration file...\n")
	for _, prober := range p.probers {
		go prober.Probe()
	}
}

func (p *Probes) registerProbes() {
	log.Printf("Registering all activated PID probes...")
	for _, probe := range p.Pid {
		p.probers = append(p.probers, &probe)
	}
	log.Printf("Registering all activated HTTP probes...")
	for _, probe := range p.Http {
		p.probers = append(p.probers, &probe)
	}
}
