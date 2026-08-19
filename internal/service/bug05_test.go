package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

func TestBug05_BatchIngestRespectsContextCancellation(t *testing.T) {
	st, svc := newTestService(t)
	incID := seedIncubator(t, st)

	readings := make([]*model.SensorReading, 50)
	for i := range readings {
		readings[i] = &model.SensorReading{
			IncubatorID: incID,
			SensorType:  model.SensorTypeTemperature,
			Value:       float64(i),
			RecordedAt:  time.Date(2024, 6, 1, 0, i, 0, 0, time.UTC),
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.Readings.BatchIngest(ctx, readings)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil (context cancellation ignored)")
	}
}
