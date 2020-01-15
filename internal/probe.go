package internal

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type Prober interface {
	Probe()
}

type Doc struct {
	PidProbes  []PidProbe  `yaml:"pid"`
	HttpProbes []HttpProbe `yaml:"http"`
}

type PidProbe struct {
	Name           string `yaml:"name"`
	Hostname       string `yaml:"hostname"`
	ServiceAccount string `yaml:"service-account"`
	Pid            int    `yaml:"pid"`
	Delay          int    `yaml:"delay"`
}

type HttpProbe struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Delay int    `yaml:"delay"`
}

func (h *HttpProbe) Probe() {
	for {
		resp, err := http.Get(h.Url)
		if err != nil {
			// handle error from call --> mark service as down...
			log.Printf("<<HTTP PROBE [%s]>> an error occured while doing HTTP call to [%s]. got '%s'\n", h.Name, h.Url, err)
			time.Sleep(time.Duration(h.Delay) * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("<<HTTP PROBE [%s]>> an error occured while reading response body. got '%s'\n", h.Name, err)
			time.Sleep(time.Duration(h.Delay) * time.Second)
			continue
		}

		log.Printf("<<HTTP PROBE [%s]>> Service returned [%s]\n", h.Name, string(body))

		resp.Body.Close()
		time.Sleep(time.Duration(h.Delay) * time.Second)
	}
}

func (p *PidProbe) Probe() {
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
