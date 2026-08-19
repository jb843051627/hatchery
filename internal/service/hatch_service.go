package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type HatchRecordService struct {
	store *store.HatchRecordStore
}

func NewHatchRecordService(s *store.HatchRecordStore) *HatchRecordService {
	return &HatchRecordService{store: s}
}

func (s *HatchRecordService) Create(ctx context.Context, r *model.HatchRecord) (int64, error) {
	if r.HealthyCount+r.WeakCount+r.DeadCount != r.HatchedCount {
		return 0, &model.ValidationError{Field: "counts", Message: "healthy + weak + dead must equal hatched"}
	}
	return s.store.Create(ctx, r)
}

func (s *HatchRecordService) Get(ctx context.Context, id int64) (*model.HatchRecord, error) {
	return s.store.GetByID(ctx, id)
}

func (s *HatchRecordService) ListByBatch(ctx context.Context, batchID int64) ([]*model.HatchRecord, error) {
	records, err := s.store.ListByBatch(ctx, batchID)
	if err != nil {
		return nil, err
	}
	out := make([]*model.HatchRecord, len(records))
	copy(out, records)
	return out, nil
}

func (s *HatchRecordService) ListByDateRange(ctx context.Context, from, to interface{}) ([]*model.HatchRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *HatchRecordService) SummaryByBatch(ctx context.Context, batchID int64) (totalHatched, healthy, weak, dead int, err error) {
	return s.store.SummaryByBatch(ctx, batchID)
}

func (s *HatchRecordService) CountByBatch(ctx context.Context, batchID int64) (int, error) {
	return s.store.CountByBatch(ctx, batchID)
}
