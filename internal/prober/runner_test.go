package prober

// mock of the interface ProbeRunner
type mockRunner struct {
	RunFn func(probe *Probe)
}

// NewMockRunner creates a new mock instance
func NewMockRunner(RunFn func(probe *Probe)) *mockRunner {
	return &mockRunner{RunFn: RunFn}
}

// Mock Run(probe *Probe) function
func (m *mockRunner) Run(probe *Probe) {
	m.RunFn(probe)
}

// TODO: how to test the run function ? :(
