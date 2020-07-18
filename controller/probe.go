// Controller contains everything that relates to HTTP(s) exposure.
// Our API is completely described here.
package controller

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/madjlzz/madprobe/internal/prober"
	"log"
	"net/http"
)

// CreateProbeRequest represents the data structure
// decoded from incoming HTTP request when trying to create a new probe.
type CreateProbeRequest struct {
	Name  string
	URL   string
	Delay uint
}

// UpdateProbeRequest represents the data structure
// decoded from incoming HTTP request when trying to update an existing probe.
type UpdateProbeRequest struct {
	Name  string
	URL   string
	Delay uint
}

// ProbeResponse represents the data structure
// send to clients when they are trying to fetch information from the API.
// It is encoded in JSON.
type ProbeResponse struct {
	Name   string
	URL    string
	Status string
	Delay  uint
}

// ProbeController is the controller
// exposing endpoints to manage probes.
type ProbeController struct {
	ProbeService prober.ProbeService
}

// NewProbeController initialize a new ProbeController
// to expose endpoints for managing probes.
func NewProbeController(ps prober.ProbeService) ProbeController {
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
		var mr *malformedContent
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = pc.ProbeService.Create(prober.Probe{
		Name:  cpr.Name,
		URL:   cpr.URL,
		Delay: cpr.Delay,
	})
	if err != nil {
		switch err {
		case prober.ErrProbeAlreadyExist:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	_, _ = fmt.Fprintf(w, "Probe [%s] has been successfuly created.", cpr.Name)
}

// Read allows consumer to retrieve a probe in the system given it's name.
// It will return a HTTP 200 status code with the probe's details if it succeeds, a human readable error otherwise.
//
// GET /api/v1/probe/{name}
func (pc *ProbeController) Read(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	probe, err := pc.ProbeService.Read(vars["name"])
	if err != nil {
		switch err {
		case prober.ErrProbeNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Encode the probe in json and send it over to the client
	pr := ProbeResponse{
		Name:   probe.Name,
		URL:    probe.URL,
		Status: probe.Status,
		Delay:  probe.Delay,
	}
	err = encodeJSONBody(w, &pr)
	if err != nil {
		var mr *malformedContent
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
}

// ReadAll allows consumer to retrieve all probe existing in the system.
// It will return a HTTP 200 status code with all probe's details if it succeeds, a human readable error otherwise.
//
// GET /api/v1/probe
func (pc *ProbeController) ReadAll(w http.ResponseWriter, req *http.Request) {
	probes := pc.ProbeService.ReadAll()

	pr := make([]ProbeResponse, 0)
	for _, value := range probes {
		pr = append(pr, ProbeResponse{
			Name:   value.Name,
			URL:    value.URL,
			Status: value.Status,
			Delay:  value.Delay,
		})
	}

	err := encodeJSONBody(w, &pr)
	if err != nil {
		var mr *malformedContent
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
}

// Update allows consumer to update an existing probe in the system.
// It will return a HTTP 200 status code if it succeeds, a human readable error otherwise.
//
// PUT /api/v1/probe/{name}
func (pc *ProbeController) Update(w http.ResponseWriter, req *http.Request) {
	var upr UpdateProbeRequest
	vars := mux.Vars(req)

	err := decodeJSONBody(w, req, &upr)
	if err != nil {
		var mr *malformedContent
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = pc.ProbeService.Update(vars["name"], prober.Probe{
		Name:  upr.Name,
		URL:   upr.URL,
		Delay: upr.Delay,
	})
	if err != nil {
		switch err {
		case prober.ErrProbeAlreadyExist:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	_, _ = fmt.Fprintf(w, "Probe [%s] has been updated.", upr.Name)
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
		case prober.ErrProbeNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	_, _ = fmt.Fprintf(w, "Probe [%s] has been successfuly deleted.", vars["name"])
}
