package prober

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/madjlzz/madprobe/internal/mock"
	"github.com/madjlzz/madprobe/internal/persistence"
	"testing"
)

func TestInsertReturnErrorOnValidationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	s := NewProbeService(nil, m)
	p := NewProbe("", "", 0)

	err := s.Insert(*p)
	if err == nil {
		t.Error("bad insert data should result in an validation error")
	}
	if _, ok := err.(*validatorError); !ok {
		t.Error("error should be a validation error")
	}
}

func TestInsertReturnErrorOnGetFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	m.
		EXPECT().
		Get(gomock.Any()).
		Return(nil, errors.New("mock Get method returns error")).
		Times(1)

	s := NewProbeService(nil, m)
	p := NewProbe("TheName", "http://localhost:8080/", 5)

	err := s.Insert(*p)
	if err == nil {
		t.Error("failing persistent layer should result in an error")
	}
}

func TestInsertReturnErrorIfProbeAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := NewProbe("TheName", "http://localhost:8080/", 5)

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	m.
		EXPECT().
		Get(gomock.Any()).
		Return(persistence.NewEntity(p.Name, p.URL, p.Delay), nil).
		Times(1)

	s := NewProbeService(nil, m)

	err := s.Insert(*p)
	if !errors.Is(err, ErrProbeAlreadyExist) {
		t.Error("returned error should be [ErrProbeAlreadyExist]")
	}
}

func TestInsertReturnErrorOnInsertFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := NewProbe("TheName", "http://localhost:8080/", 5)

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)
	m.EXPECT().Get(gomock.Any()).Times(1)

	m.
		EXPECT().
		Insert(gomock.Eq(persistence.NewEntity(p.Name, p.URL, p.Delay))).
		Return(errors.New("mock Insert method returns error")).
		Times(1)

	s := NewProbeService(nil, m)

	err := s.Insert(*p)
	if err == nil {
		t.Error("failing persistent layer should result in an error")
	}
}

func TestInsertSuccessUpdateCacheAndReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := NewProbe("TheName", "http://localhost:8080/", 5)

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)
	m.EXPECT().Get(gomock.Any()).Times(1)

	mockRunner := NewMockRunner(func(probe *Probe) {
		fmt.Println("[MOCK] Running probes...")
	})

	m.
		EXPECT().
		Insert(gomock.Eq(persistence.NewEntity(p.Name, p.URL, p.Delay))).
		Times(1)

	s := NewProbeService(mockRunner, m)

	_ = s.Insert(*p)
	if _, ok := s.probes[p.Name]; !ok {
		t.Error("internal service cache should be updated to keep the state of running probes!")
	}
}

func TestGetReturnErrorOnGetFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	m.EXPECT().
		Get(gomock.Any()).
		Return(nil, errors.New("mock Get method returns error")).
		Times(1)

	s := NewProbeService(nil, m)

	_, err := s.Get("TheName")
	if err == nil {
		t.Error("failing persistent layer should result in an error")
	}
}

func TestGetReturnErrProbeNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	m.EXPECT().
		Get(gomock.Any()).
		Return(nil, nil).
		Times(1)

	s := NewProbeService(nil, m)

	_, err := s.Get("TheName")
	if !errors.Is(err, ErrProbeNotFound) {
		t.Error("returned error should be [ErrProbeNotFound]")
	}
}

func TestGetSuccessReturnProbeFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entity := persistence.NewEntity("TheName", "http://localhost/", 5)

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	m.EXPECT().
		Get(gomock.Any()).
		Return(entity, nil).
		Times(1)

	s := NewProbeService(nil, m)
	s.probes[entity.Name] = NewProbe(entity.Name, entity.URL, entity.Delay)

	probe, err := s.Get("TheName")
	if err != nil {
		t.Error("no error should have been registered")
	}
	if probe == nil {
		t.Error("a probe should have been retrieved from the cache!")
	}
}

func TestGetAllReturnErrorOnGetAllFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	firstCall := m.EXPECT().GetAll()
	secondCall := m.
		EXPECT().
		GetAll().
		Return(nil, errors.New("mock GetAll method returns error"))

	gomock.InOrder(firstCall, secondCall)

	s := NewProbeService(nil, m)

	_, err := s.GetAll()
	if err == nil {
		t.Error("failing persistent layer should result in an error")
	}
}

func TestGetAllReturnProbesFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entities := []*persistence.Entity{
		persistence.NewEntity("TheName", "http://localhost/", 5),
	}

	m := mock.NewMockPersister(ctrl)

	firstCall := m.EXPECT().GetAll()
	secondCall := m.
		EXPECT().
		GetAll().
		Return(entities, nil)

	gomock.InOrder(firstCall, secondCall)

	s := NewProbeService(nil, m)
	s.probes[entities[0].Name] = NewProbe(entities[0].Name, entities[0].URL, entities[0].Delay)

	probes, err := s.GetAll()
	if err != nil {
		t.Error("no error should have been registered")
	}
	if len(probes) != 1 {
		t.Error("all probes should have been retrieved from the cache!")
	}
}

func TestDeleteReturnErrorOnValidationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	s := NewProbeService(nil, m)
	probeName := ""

	err := s.Delete(probeName)
	if err == nil {
		t.Error("bad insert data should result in an validation error")
	}
	if _, ok := err.(*validatorError); !ok {
		t.Error("error should be a validation error")
	}
}

func TestDeleteReturnErrorOnDeleteFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	m.
		EXPECT().
		Delete(gomock.Any()).
		Return(errors.New("mock Delete method returns error")).
		Times(1)

	s := NewProbeService(nil, m)
	probeName := "TheProbe"

	err := s.Delete(probeName)
	if err == nil {
		t.Error("failing persistent layer should result in an error")
	}
}

// TODO: test blocking because of the channel. Why is that ?
/*func TestDeleteSuccessUpdateCacheAndProbeState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Times(1)

	m.
		EXPECT().
		Delete(gomock.Any()).
		Return(nil).
		Times(1)

	s := NewProbeService(nil, m)
	probeName := "TheProbe"

	// Inserting fake cache data just to check if Delete is updating properly the cache.
	s.probes[probeName] = NewProbe(probeName, "TheURL", 5)
	channelRef := s.probes[probeName].Finish

	_ = s.Delete(probeName)
	if !<-channelRef {
		t.Error("channel should be updated to manage the probe state.")
	}
	if len(s.probes) != 0 {
		t.Error("cache should have been updated after the deletion of the probe.")
	}

}*/
