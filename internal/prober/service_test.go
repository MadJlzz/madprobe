package prober

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/madjlzz/madprobe/internal/mock"
	"github.com/madjlzz/madprobe/internal/persistence"
	"testing"
)

func TestInsertReturnErrorOnValidationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockPersister(ctrl)
	m.EXPECT().GetAll().Return([]*persistence.Entity{}, nil).AnyTimes()

	s := NewProbeService(nil, m, nil)
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
	m.EXPECT().GetAll().Return([]*persistence.Entity{}, nil).AnyTimes()

	m.
		EXPECT().
		Get(gomock.Any()).
		Return(nil, errors.New("mock Get method returns error")).
		Times(1)

	s := NewProbeService(nil, m, nil)
	p := NewProbe("TheName", "http://localhost:8080/", 5)

	err := s.Insert(*p)
	if err == nil {
		t.Error("failing persistent layer should result in an error")
	}
}