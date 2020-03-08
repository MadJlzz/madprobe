// Controller contains everything that relates to HTTP(s) exposure.
// Our API is completely described here.
package controller

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/madjlzz/madprobe/internal/service"
	"log"
	"net/http"
)

// ProbeService represent the interface used
// to manipulate probes.
type ProbeService interface {
	Create(probe service.Probe) error
	Delete(name string) error
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

	err = pc.ProbeService.Create(service.Probe{
		Name:  cpr.Name,
		URL:   cpr.URL,
		Delay: cpr.Delay,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, "Probe [%s] has been successfuly created.", cpr.Name)
}

// Delete allows consumer to delete an existing probe in the system.
// It will return a HTTP 200 status code if it succeeds, a human readable error otherwise.
//
// DELETE /api/v1/probe/{name}
func (pc *ProbeController) Delete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	err := pc.ProbeService.Delete(vars["name"])
	if err != nil {
		switch err {
		case service.ErrProbeNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	_, _ = fmt.Fprintf(w, "Probe [%s] has been successfuly deleted.", vars["name"])
}
