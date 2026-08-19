package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

func TestBug07_BatchCreateRollbackOnError(t *testing.T) {
	st, svc := newTestService(t)
	incID := seedIncubator(t, st)

	tasks := []*model.MaintenanceTask{
		{IncubatorID: incID, Description: "task-1", ScheduledFor: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)},
		{IncubatorID: -1, Description: "bad-inc", ScheduledFor: time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC)},
		{IncubatorID: incID, Description: "task-3", ScheduledFor: time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC)},
	}

	err := svc.Maintenance.BatchCreate(context.Background(), tasks)
	if err == nil {
		t.Fatal("expected error from batch with invalid incubator_id, got nil")
	}

	pending, err := svc.Maintenance.ListPending(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range pending {
		if m.Description == "task-1" || m.Description == "task-3" {
			t.Fatalf("found residual task %q after failed batch (expected rollback)", m.Description)
		}
	}
}
