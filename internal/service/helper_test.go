package service

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jb843051627/hatchery/internal/cache"
	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

func newTestService(t *testing.T) (*store.Store, *Service) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	st, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	rc := cache.NewReadingCache()
	return st, NewService(st, rc)
}

func seedIncubator(t *testing.T, st *store.Store) int64 {
	t.Helper()
	id, err := st.Incubators.Create(context.Background(), "test-inc", "loc-A", 500, model.IncubatorStatusActive, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("create incubator: %v", err)
	}
	return id
}

func seedBatch(t *testing.T, st *store.Store, incID int64, species string) int64 {
	t.Helper()
	id, err := st.Batches.Create(context.Background(), incID, 100, species,
		time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 5, 22, 0, 0, 0, 0, time.UTC),
		model.BatchStatusIncubating)
	if err != nil {
		t.Fatalf("create batch: %v", err)
	}
	return id
}

func seedEggTray(t *testing.T, st *store.Store, batchID int64, trayNum, eggCount int) {
	t.Helper()
	_, err := st.EggTrays.Create(context.Background(), batchID, trayNum, eggCount, 1.5, "farm-A")
	if err != nil {
		t.Fatalf("create egg tray: %v", err)
	}
}

func seedHatchRecord(t *testing.T, st *store.Store, batchID int64, hatched, healthy, weak, dead int, date time.Time) {
	t.Helper()
	r := &model.HatchRecord{
		BatchID:      batchID,
		HatchedCount:  hatched,
		HealthyCount:  healthy,
		WeakCount:     weak,
		DeadCount:     dead,
		HatchDate:     date,
		GradedBy:      "tester",
	}
	_, err := st.HatchRecords.Create(context.Background(), r)
	if err != nil {
		t.Fatalf("create hatch record: %v", err)
	}
}
