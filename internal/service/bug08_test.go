package service

import (
	"context"
	"sort"
	"testing"
	"time"
)

func TestBug08_ListHatchRecordsReturnsCopyNotCacheRef(t *testing.T) {
	st, svc := newTestService(t)
	incID := seedIncubator(t, st)
	batchID := seedBatch(t, st, incID, "species-B")
	seedHatchRecord(t, st, batchID, 10, 8, 1, 1, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))
	seedHatchRecord(t, st, batchID, 12, 10, 1, 1, time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC))
	seedHatchRecord(t, st, batchID, 8, 6, 1, 1, time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC))

	ctx := context.Background()
	first, err := svc.HatchRecords.ListByBatch(ctx, batchID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 records, got %d", len(first))
	}
	sort.Slice(first, func(i, j int) bool {
		return first[i].HatchedCount > first[j].HatchedCount
	})

	second, err := svc.HatchRecords.ListByBatch(ctx, batchID)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(second); i++ {
		if second[i].HatchDate.Before(second[i-1].HatchDate) {
			t.Fatalf("cache polluted: expected hatch_date ascending, got [%d] before [%d]",
				i, i-1)
		}
	}
}
