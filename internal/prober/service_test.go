package prober

import (
	"fmt"
	"github.com/madjlzz/madprobe/internal/persistence"
	"testing"
)

type mockPersistentClient struct{}

func (p *mockPersistentClient) Insert(entity *persistence.Entity) error {
	fmt.Println("[INFO] Insert mock.")
	return nil
}

func (p *mockPersistentClient) Delete(name string) error {
	fmt.Println("[INFO] Delete mock.")
	return nil
}

func (p *mockPersistentClient) Get(name string) (*persistence.Entity, error) {
	fmt.Println("[INFO] Get mock.")
	return nil, nil
}

func (p *mockPersistentClient) GetAll() ([]*persistence.Entity, error) {
	fmt.Println("[INFO] GetAll mock.")
	return nil, nil
}

func TestNewProbeServiceIsSingleton(t *testing.T) {
	s1, s2 := NewProbeService(nil, nil, nil), NewProbeService(nil, nil, nil)
	if s1 != s2 {
		t.Errorf("service should be a singleton. got [%+v], want [%+v]\n", s1, s2)
	}
}

func TestCreateReturnErrorOnValidationFailure(t *testing.T) {
	// s := NewProbeService(nil, &mockPersistentClient{}, nil)
	// p := NewProbe("", "", 0)
}
