package service

import (
	"context"
	"testing"
	"time"
)

func TestBug10_ExportReadingsPreservesOriginalTime(t *testing.T) {
	st, svc := newTestService(t)
	incID := seedIncubator(t, st)
	ctx := context.Background()

	recordedAt := time.Date(2024, 6, 15, 10, 30, 0, 0, time.FixedZone("CST", 8*3600))
	_, err := st.Readings.Create(ctx, incID, "temperature", 37.5, recordedAt)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := svc.Report.ExportReadings(ctx, incID,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if row.RecordedAt == "" {
			continue
		}
		parsed, err := time.Parse("2006-01-02 15:04:05", row.RecordedAt)
		if err != nil {
			t.Fatalf("parse recorded_at %q: %v", row.RecordedAt, err)
		}
		if parsed.Hour() != 10 {
			t.Fatalf("expected hour=10 (CST), got hour=%d in %q (UTC offset applied)",
				parsed.Hour(), row.RecordedAt)
		}
	}
}
