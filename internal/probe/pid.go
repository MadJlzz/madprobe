package probe

import (
	"log"
	"os/exec"
	"strconv"
	"time"
)

type Pid struct {
	Name           string `yaml:"name"`
	Hostname       string `yaml:"hostname"`
	ServiceAccount string `yaml:"service-account"`
	Pid            int    `yaml:"pid"`
	Delay          int    `yaml:"delay"`
}

func (p *Pid) Probe() {
	if p.Hostname != "localhost" {
		log.Printf("<<PID PROBE>> PID probe not is currently not implemented for remote check.")
		return
	}
	for {
		pid := strconv.Itoa(p.Pid)
		cmd := exec.Command("ps", "-p", pid)
		if err := cmd.Run(); err != nil {
			// handle error from call --> mark service as down...
			log.Printf("<<PID PROBE>> Process with PID [%d] not found. got '%s'\n", p.Pid, err)
			time.Sleep(time.Duration(p.Delay) * time.Second)
			continue
		}
		log.Printf("<<PID PROBE>> Process with PID [%d] is currently running.", p.Pid)
		time.Sleep(time.Duration(p.Delay) * time.Second)
	}
}
