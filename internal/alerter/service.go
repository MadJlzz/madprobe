package alerter

import (
	"errors"
	"fmt"
	"github.com/madjlzz/madprobe/internal/prober"
	"sync"
)

// Error thrown whenever the alert bus passed to the service is not initialized.
var ErrAlertBusNotReady = errors.New("bus should be initialized for alerting to work")

// Once is an object that will perform exactly one action.
// It is used to ensure the service is a singleton.
var once sync.Once
var instance *service

type service struct {
	alertBus <-chan prober.Probe
	alerters []Alerter
}

// Initialize the alerting service with existing implementations.
func NewService(alertBus <-chan prober.Probe) (*service, error) {
	if alertBus == nil {
		return nil, ErrAlertBusNotReady
	}
	alerters := []Alerter{NewDiscordAlerter()}
	once.Do(func() {
		instance = &service{
			alertBus: alertBus,
			alerters: filter(alerters),
		}
	})
	return instance, nil
}

// Run every alerter that has been correctly instantiated.
func (s *service) Run() {
	for _, a := range s.alerters {
		go a.Alert(s.alertBus)
	}
}

// Close every alerter that can be closed.
func (s *service) Close() error {
	var err error
	for _, a := range s.alerters {
		switch t := a.(type) {
		case AlertCloser:
			cErr := t.Close()
			if cErr != nil {
				err = fmt.Errorf("%w\n", cErr)
			}
		}
	}
	return err
}

func filter(alerters []Alerter) []Alerter {
	var fAlerters []Alerter
	for _, a := range alerters {
		if a != nil {
			fAlerters = append(fAlerters, a)
		}
	}
	return fAlerters
}