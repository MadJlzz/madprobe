// Controller contains everything that relates to HTTP(s) exposure.
// Our API is completely described here.
package controller

import (
	"errors"
	"fmt"
	"github.com/madjlzz/madprobe/internal/service"
	"log"
	"net/http"
)

// ProbeService represent the interface used
// to manipulate probes.
type ProbeService interface {
	Create(probe service.Probe) error
}

// CreateProbeRequest represents the data structure
// decoded from incoming HTTP request when trying to create a new probe.
type CreateProbeRequest struct {
	Name  string
	URL   string
	Delay uint
}

// ProbeController is the controller
// exposing endpoints to manage probes.
type ProbeController struct {
	ProbeService ProbeService
}

// NewProbeController initialize a new ProbeController
// to expose endpoints for managing probes.
func NewProbeController(ps ProbeService) ProbeController {
	return ProbeController{
		ProbeService: ps,
	}
}

// Create allows consumer to create a new probe in the system.
// It will return a HTTP 200 status code if it succeeds, a human readable error otherwise.
//
// POST /api/v1/probe/create
func (pc *ProbeController) Create(w http.ResponseWriter, req *http.Request) {
	var cpr CreateProbeRequest

	err := decodeJSONBody(w, req, &cpr)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Call the service layer in order to persist the probe.
	// Bad thing is that we cannot mock the call as is.
	err = pc.ProbeService.Create(service.Probe{
		Name:  cpr.Name,
		URL:   cpr.URL,
		Delay: cpr.Delay,
	})

	_, _ = fmt.Fprintf(w, "CreateProbeRequest: %+v", cpr)
}
