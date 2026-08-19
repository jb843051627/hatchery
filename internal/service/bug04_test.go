package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/hatchery/internal/model"
)

func TestBug04_AlertCreateErrorWrapsCorrectly(t *testing.T) {
	st, svc := newTestService(t)
	incID := seedIncubator(t, st)
	ctx := context.Background()

	_, err := svc.Alerts.Create(ctx, incID, 0, "critical", "temp too high")
	if err != nil {
		t.Fatalf("create alert: %v", err)
	}

	_, err = svc.Alerts.Create(ctx, incID, 0, "critical", "")
	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected errors.As(err, *ValidationError), got: %v (type %T)", err, err)
	}
}
