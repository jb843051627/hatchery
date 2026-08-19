package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type IncubatorService struct {
	store *store.IncubatorStore
}

func NewIncubatorService(s *store.IncubatorStore) *IncubatorService {
	return &IncubatorService{store: s}
}

func (s *IncubatorService) Get(ctx context.Context, id int64) (*model.Incubator, error) {
	inc, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get incubator: %w", err)
	}
	return inc, nil
}

func (s *IncubatorService) Create(ctx context.Context, name, location string, capacity int, status string) (int64, error) {
	if status == "" {
		status = model.IncubatorStatusActive
	}
	return s.store.Create(ctx, name, location, capacity, status, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
}

func (s *IncubatorService) ListAll(ctx context.Context) ([]*model.Incubator, error) {
	return s.store.ListAll(ctx)
}

func (s *IncubatorService) ListByStatus(ctx context.Context, status string) ([]*model.Incubator, error) {
	return s.store.ListByStatus(ctx, status)
}

func (s *IncubatorService) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := s.store.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get incubator: %w", err)
	}
	return s.store.UpdateStatus(ctx, id, status)
}

func (s *IncubatorService) Count(ctx context.Context) (int, error) {
	return s.store.Count(ctx)
}
