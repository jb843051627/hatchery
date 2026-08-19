package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/hatchery/internal/model"
)

func TestBug01_GetNonExistentIncubatorReturnsError(t *testing.T) {
	_, svc := newTestService(t)
	inc, err := svc.Incubators.Get(context.Background(), 999999)
	if err == nil {
		t.Fatalf("expected error for non-existent incubator, got inc=%v err=nil", inc)
	}
	if !errors.Is(err, model.ErrIncubatorNotFound) {
		t.Fatalf("expected ErrIncubatorNotFound, got: %v", err)
	}
}
