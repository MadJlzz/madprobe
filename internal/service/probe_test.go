package service

import (
	"testing"
)

func TestNewProbeService(t *testing.T) {
	s1, s2 := NewProbeService(nil, nil), NewProbeService(nil, nil)
	if s1 != s2 {
		t.Errorf("ProbeService should be a singleton. got %+v, want %+v\n", s1, s2)
	}
}
