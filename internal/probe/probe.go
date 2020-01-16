package probe

import (
	"github.com/madjlzz/madprobe/internal/server"
	"log"
)

const (
	StatusUp   = "UP"
	StatusDown = "DOWN"
)

type Probes struct {
	Pid  []Pid  `yaml:"pid"`
	Http []Http `yaml:"http"`
}

type Probe struct {
	Name  string `yaml:"name"`
	Delay int    `yaml:"delay"`
	app   *server.App
}

func Run(p *Probes, app *server.App) {
	log.Printf("Running probes defined in yaml configuration file...\n")
	for _, probe := range p.Pid {
		probe.app = app
		go probe.Check()
	}
	for _, probe := range p.Http {
		probe.app = app
		go probe.Check()
	}
}
