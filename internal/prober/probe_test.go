package prober

import "testing"

func TestNewProbe(t *testing.T) {
	probe := NewProbe("TheProbe", "TheURL", 5)
	if probe.Name != "TheProbe" {
		t.Errorf("Name property should be [TheProbe]. got: %s\n", probe.Name)
	}
	if probe.URL != "TheURL" {
		t.Errorf("URL property should be [TheURL]. got: %s\n", probe.URL)
	}
	if probe.Delay != 5 {
		t.Errorf("Delay property should be [5]. got: %d\n", probe.Delay)
	}
	if probe.Status != "" {
		t.Errorf("Status property should be [\"\"]. got: %s\n", probe.Status)
	}
	if probe.Finish == nil {
		t.Errorf("Finish channel should be initialized. got %v\n", probe.Finish)
	}
}
