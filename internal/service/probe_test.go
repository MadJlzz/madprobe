package service

import (
	"testing"
)

func TestNewProbeService(t *testing.T) {
	s1, _ := NewProbeService()
	s2, _ := NewProbeService()
	if s1 != s2 {
		t.Errorf("ProbeService should be a singleton. got %d, want %d\n", s1, s2)
	}
}
