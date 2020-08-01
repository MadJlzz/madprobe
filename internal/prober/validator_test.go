package prober

import (
	"testing"
)

func TestError(t *testing.T) {
	err := validatorError{
		field: "name",
		msg:   "name should not be empty",
	}
	errAsStr := err.Error()
	if errAsStr != "Field [name]: name should not be empty\n" {
		t.Errorf("the returned error string is invalid: [Field [name]: name should not be empty]. got %s", errAsStr)
	}
}

func TestNameInvalid(t *testing.T) {
	probe := NewProbe("", "", 0)
	err := nameInvalid(*probe)
	if err == nil {
		t.Errorf("the probe's name is invalid. an error should be returned.")
	}
	if e, ok := err.(*validatorError); ok {
		if e.field != "Name" {
			t.Errorf("validatorError field must be [Name]. got: %s\n", e.field)
		}
		if e.msg != "name is required" {
			t.Errorf("validatorError msg must be [name is required]. got: %s\n", e.msg)
		}
	}
}

func TestNameValid(t *testing.T) {
	probe := NewProbe("ValidName", "", 0)
	err := nameInvalid(*probe)
	if err != nil {
		t.Errorf("no error should be thrown with a valid probe's name.")
	}
}

func TestEmptyURLInvalid(t *testing.T) {
	probe := NewProbe("", "", 0)
	err := urlInvalid(*probe)
	if err == nil {
		t.Errorf("the probe's URL is invalid. an error should be returned.")
	}
	if e, ok := err.(*validatorError); ok {
		if e.field != "URL" {
			t.Errorf("validatorError field must be [URL]. got: %s\n", e.field)
		}
		if e.msg != "URL is required" {
			t.Errorf("validatorError msg must be [URL is required]. got: %s\n", e.msg)
		}
	}
}

func TestMalformedURLInvalid(t *testing.T) {
	probe := NewProbe("", "localhost:8080", 0)
	err := urlInvalid(*probe)
	if err == nil {
		t.Errorf("the probe's URL is invalid. an error should be returned.")
	}
	if e, ok := err.(*validatorError); ok {
		if e.field != "URL" {
			t.Errorf("validatorError field must be [URL]. got: %s\n", e.field)
		}
		if e.msg != "URL is malformed" {
			t.Errorf("validatorError msg must be [URL is malformed]. got: %s\n", e.msg)
		}
	}
}

func TestURLValid(t *testing.T) {
	probe := NewProbe("", "http://localhost/", 0)
	err := urlInvalid(*probe)
	if err != nil {
		t.Errorf("no error should be thrown with a valid probe's URL.")
	}
}

func TestDelayInvalid(t *testing.T) {
	probe := NewProbe("", "", 0)
	err := delayInvalid(*probe)
	if err == nil {
		t.Errorf("the probe's delay is invalid. an error should be returned.")
	}
	if e, ok := err.(*validatorError); ok {
		if e.field != "Delay" {
			t.Errorf("validatorError field must be [Delay]. got: %s\n", e.field)
		}
		if e.msg != "Delay must be at least 1 and strictly positive" {
			t.Errorf("validatorError msg must be [Delay must be at least 1 and strictly positive]. got: %s\n", e.msg)
		}
	}
}

func TestDelayValid(t *testing.T) {
	probe := NewProbe("ValidName", "", 1)
	err := delayInvalid(*probe)
	if err != nil {
		t.Errorf("no error should be thrown with a valid probe's delay.")
	}
}

func TestRunValidatorWithError(t *testing.T) {
	probe := NewProbe("ValidName", "", 1)
	err := runValidators(*probe, nameInvalid, urlInvalid, delayInvalid)
	if err == nil {
		t.Errorf("the probe's URL is invalid. an error should be returned.")
	}
	if e, ok := err.(*validatorError); ok {
		if e.field != "URL" {
			t.Errorf("validatorError field must be [URL]. got: %s\n", e.field)
		}
		if e.msg != "URL is required" {
			t.Errorf("validatorError msg must be [URL is required]. got: %s\n", e.msg)
		}
	}
}

func TestRunValidatorWithNoError(t *testing.T) {
	probe := NewProbe("ValidName", "http://localhost/", 1)
	err := runValidators(*probe, nameInvalid, urlInvalid, delayInvalid)
	if err != nil {
		t.Errorf("no error should be thrown with a valid probe.")
	}
}
