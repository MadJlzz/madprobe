package prober

import (
	"fmt"
	"net/url"
)

// Custom error type that occurs when there is a validation error.
type validatorError struct {
	field string
	msg   string
}

// Implementation of the error interface.
// Prints out a validationError.
func (ve *validatorError) Error() string {
	return fmt.Sprintf("Field [%s]: %s\n", ve.field, ve.msg)
}

// Validate the name property of the probe.
// Returns an error if the name is empty.
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

// Validate the delay property of the probe.
// Returns an error if the delay is 0 or negative.
func delayInvalid(probe Probe) error {
	if probe.Delay <= 0 {
		return &validatorError{
			field: "Delay",
			msg:   "Delay must be at least 1 and strictly positive",
		}
	}
	return nil
}

// Handy type that allow us to pass a function that takes a probe
// for validation.
type validateFunc func(probe Probe) error

// Runs all of the given validator functions for the passed probe.
func runValidators(probe Probe, fns ...validateFunc) error {
	for _, fn := range fns {
		if err := fn(probe); err != nil {
			return err
		}
	}
	return nil
}
