package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

func TestBug09_CreateRejectsInvalidCounts(t *testing.T) {
	st, svc := newTestService(t)
	incID := seedIncubator(t, st)
	batchID := seedBatch(t, st, incID, "species-C")

	r := &model.HatchRecord{
		BatchID:      batchID,
		HatchedCount:  10,
		HealthyCount: 5,
		WeakCount:    3,
		DeadCount:    1,
		HatchDate:    time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		GradedBy:     "tester",
	}
	_, err := svc.HatchRecords.Create(context.Background(), r)
	if err == nil {
		t.Fatal("expected validation error (5+3+1=9 != 10), got nil")
	}
}
