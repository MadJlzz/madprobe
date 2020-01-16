package probe

import (
	"log"
	"os/exec"
	"strconv"
	"time"
)

type Pid struct {
	Probe          `yaml:",inline"`
	Hostname       string `yaml:"hostname"`
	ServiceAccount string `yaml:"service-account"`
	Pid            int    `yaml:"pid"`
}

func (p *Pid) Check() {
	if p.Hostname != "localhost" {
		log.Printf("<<PID PROBE>> PID probe not is currently not implemented for remote check.")
		return
	}
	for {
		pid := strconv.Itoa(p.Pid)
		cmd := exec.Command("ps", "-p", pid)
		if err := cmd.Run(); err != nil {
			log.Printf("<<PID PROBE>> Process with PID [%d] not found. got '%s'\n", p.Pid, err)
			p.app.UpdateTemplateData(p.Name, StatusDown)
		} else {
			log.Printf("<<PID PROBE>> Process with PID [%d] is currently running.", p.Pid)
			p.app.UpdateTemplateData(p.Name, StatusUp)
		}
		time.Sleep(time.Duration(p.Delay) * time.Second)
	}
}
