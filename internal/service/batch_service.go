package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type BatchService struct {
	batches    *store.BatchStore
	incubators *store.IncubatorStore
}

func NewBatchService(b *store.BatchStore, i *store.IncubatorStore) *BatchService {
	return &BatchService{batches: b, incubators: i}
}

func (s *BatchService) Create(ctx context.Context, incubatorID int64, eggCount int, species string, startDate, expectedHatch time.Time) (int64, error) {
	inc, err := s.incubators.GetByID(ctx, incubatorID)
	if err != nil {
		return 0, fmt.Errorf("get incubator: %w", err)
	}
	if inc == nil {
		return 0, model.ErrIncubatorNotFound
	}
	if inc.Status != model.IncubatorStatusActive {
		return 0, &model.ValidationError{Field: "incubator_status", Message: "incubator is not active"}
	}
	return s.batches.Create(ctx, incubatorID, eggCount, species, startDate, expectedHatch, model.BatchStatusPending)
}

func (s *BatchService) Get(ctx context.Context, id int64) (*model.Batch, error) {
	b, err := s.batches.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get batch: %w", err)
	}
	return b, nil
}

func (s *BatchService) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.Batch, error) {
	return s.batches.ListByIncubator(ctx, incubatorID)
}

func (s *BatchService) UpdateStatus(ctx context.Context, id int64, status string) error {
	b, err := s.batches.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get batch: %w", err)
	}
	if b == nil {
		return model.ErrBatchNotFound
	}
	if b.Status == model.BatchStatusCompleted && status != model.BatchStatusCompleted {
		return &model.ValidationError{Field: "status", Message: "cannot change completed batch status"}
	}
	return s.batches.UpdateStatus(ctx, id, status)
}

func (s *BatchService) ListByStatus(ctx context.Context, status string) ([]*model.Batch, error) {
	return s.batches.ListByStatus(ctx, status)
}

func (s *BatchService) StartIncubation(ctx context.Context, id int64) error {
	b, err := s.batches.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get batch: %w", err)
	}
	if b == nil {
		return model.ErrBatchNotFound
	}
	if b.Status != model.BatchStatusPending {
		return &model.ValidationError{Field: "status", Message: "batch must be pending to start incubation"}
	}
	return s.batches.UpdateStatus(ctx, id, model.BatchStatusIncubating)
}

func (s *BatchService) GetActiveBatch(ctx context.Context, incubatorID int64) (*model.Batch, error) {
	return s.batches.GetActiveBatch(ctx, incubatorID)
}

func (s *BatchService) CountByStatus(ctx context.Context, status string) (int, error) {
	return s.batches.CountByStatus(ctx, status)
}
