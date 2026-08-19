package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/hatchery/internal/model"
)

func TestBug06_GetNonExistentScheduleReturnsError(t *testing.T) {
	_, svc := newTestService(t)
	sc, err := svc.Schedules.Get(context.Background(), 999999)
	if err == nil {
		t.Fatalf("expected error for non-existent schedule, got sc=%v err=nil", sc)
	}
	if !errors.Is(err, model.ErrScheduleNotFound) {
		t.Fatalf("expected ErrScheduleNotFound, got: %v", err)
	}
}
