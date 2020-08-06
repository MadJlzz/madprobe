package prober

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// runner is an implementation of ProbeRunner
type runner struct {
	client   *http.Client
	alertBus chan<- Probe
}

// NewProbeRunner allow to create a new probe runner.
func NewProbeRunner(httpClient *http.Client, alertBus chan<- Probe) *runner {
	return &runner{
		client:   httpClient,
		alertBus: alertBus,
	}
}

// run launches probes in a separate goroutine.
func (r *runner) Run(probe *Probe) {
	var oldStatus string
	for {
		select {
		case <-probe.Finish:
			log.Printf("<<HTTP PROBE [%s]>> Stopping probe...\n", probe.Name)
			return
		default:
			resp, err := r.client.Get(probe.URL)
			oldStatus = probe.Status
			if err != nil {
				probe.Status = downStatus
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] is down.\n", probe.Name, probe.URL)
			} else if resp.StatusCode != 200 {
				b, _ := ioutil.ReadAll(resp.Body)
				probe.Status = downStatus
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] returned an error. got: ['%v']\n", probe.Name, probe.URL, string(b))
			} else {
				probe.Status = upStatus
				log.Printf("<<HTTP(s) PROBE [%s]>> Service targeting [%s] is alive.\n", probe.Name, probe.URL)
				_ = resp.Body.Close()
			}
		}
		// If the status has changed, we can send an event to the alerter bus...
		if oldStatus != probe.Status {
			r.alertBus <- *probe
		}
		time.Sleep(time.Duration(probe.Delay) * time.Second)
	}
}
