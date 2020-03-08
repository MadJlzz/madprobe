package service

import (
	"fmt"
	"net/url"
)

type validatorError struct {
	field string
	msg   string
}

func (ve *validatorError) Error() string {
	return fmt.Sprintf("Field [%s]: %s\n", ve.field, ve.msg)
}

type validateFunc func(probe Probe) error

func nameInvalid(probe Probe) error {
	if probe.Name == "" {
		return &validatorError{
			field: "Name",
			msg:   "name is required",
		}
	}
	return nil
}

func urlInvalid(probe Probe) error {
	if probe.URL == "" {
		return &validatorError{
			field: "URL",
			msg:   "URL is required",
		}
	}
	_, err := url.ParseRequestURI(probe.URL)
	if err != nil {
		return &validatorError{
			field: "URL",
			msg:   "URL is malformed",
		}
	}
	u, err := url.Parse(probe.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return &validatorError{
			field: "URL",
			msg:   "URL is malformed",
		}
	}
	return nil
}

func delayInvalid(probe Probe) error {
	if probe.Delay <= 0 {
		return &validatorError{
			field: "Delay",
			msg:   "Delay must be at least 1 and strictly positive",
		}
	}
	return nil
}

func runValidators(probe Probe, fns ...validateFunc) error {
	for _, fn := range fns {
		if err := fn(probe); err != nil {
			return err
		}
	}
	return nil
}
