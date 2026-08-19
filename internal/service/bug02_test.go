package service

import (
	"context"
	"sort"
	"testing"
)

func TestBug02_ListByBatchReturnsCopyNotCacheRef(t *testing.T) {
	st, svc := newTestService(t)
	incID := seedIncubator(t, st)
	batchID := seedBatch(t, st, incID, "species-A")
	seedEggTray(t, st, batchID, 1, 30)
	seedEggTray(t, st, batchID, 2, 24)
	seedEggTray(t, st, batchID, 3, 28)

	ctx := context.Background()
	first, err := svc.EggTrays.ListByBatch(ctx, batchID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 trays, got %d", len(first))
	}
	sort.Slice(first, func(i, j int) bool {
		return first[i].TrayNumber > first[j].TrayNumber
	})

	second, err := svc.EggTrays.ListByBatch(ctx, batchID)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(second); i++ {
		if second[i].TrayNumber < second[i-1].TrayNumber {
			t.Fatalf("cache polluted: expected tray_number ascending, got [%d]=%d < [%d]=%d",
				i, second[i].TrayNumber, i-1, second[i-1].TrayNumber)
		}
	}
}
