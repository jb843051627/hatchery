package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type ScheduleService struct {
	store *store.ScheduleStore
}

func NewScheduleService(s *store.ScheduleStore) *ScheduleService {
	return &ScheduleService{store: s}
}

func (s *ScheduleService) Create(ctx context.Context, sc *model.Schedule) (int64, error) {
	if sc.Status == "" {
		sc.Status = model.ScheduleStatusPending
	}
	conflict, err := s.store.FindConflict(ctx, sc.IncubatorID, sc.StartTime, sc.EndTime)
	if err != nil {
		return 0, fmt.Errorf("check conflict: %w", err)
	}
	if conflict != nil {
		return 0, model.ErrScheduleConflict
	}
	return s.store.Create(ctx, sc)
}

func (s *ScheduleService) Get(ctx context.Context, id int64) (*model.Schedule, error) {
	sc, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	if sc == nil {
		return nil, model.ErrScheduleNotFound
	}
	return sc, nil
}

func (s *ScheduleService) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.Schedule, error) {
	return s.store.ListByIncubator(ctx, incubatorID)
}

func (s *ScheduleService) ListByBatch(ctx context.Context, batchID int64) ([]*model.Schedule, error) {
	return s.store.ListByBatch(ctx, batchID)
}

func (s *ScheduleService) Activate(ctx context.Context, id int64) error {
	sc, err := s.store.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get schedule: %w", err)
	}
	if sc == nil {
		return model.ErrScheduleNotFound
	}
	if sc.Status != model.ScheduleStatusPending {
		return &model.ValidationError{Field: "status", Message: "schedule must be pending to activate"}
	}
	return s.store.UpdateStatus(ctx, id, model.ScheduleStatusActive)
}

func (s *ScheduleService) Cancel(ctx context.Context, id int64) error {
	return s.store.UpdateStatus(ctx, id, model.ScheduleStatusCancelled)
}

func (s *ScheduleService) ListActive(ctx context.Context) ([]*model.Schedule, error) {
	return s.store.ListActive(ctx)
}
