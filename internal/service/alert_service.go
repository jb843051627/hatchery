package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type AlertService struct {
	store *store.AlertStore
}

func NewAlertService(s *store.AlertStore) *AlertService {
	return &AlertService{store: s}
}

func (s *AlertService) Create(ctx context.Context, incubatorID, batchID int64, level, message string) (int64, error) {
	if !model.ValidAlertLevel(level) {
		return 0, &model.ValidationError{Field: "level", Message: "invalid alert level"}
	}
	a := &model.Alert{
		IncubatorID: incubatorID,
		BatchID:     batchID,
		Level:       level,
		Message:     message,
		Status:      model.AlertStatusActive,
	}
	id, err := s.store.Create(ctx, a)
	if err != nil {
		return 0, fmt.Errorf("create alert: %v", err)
	}
	return id, nil
}

func (s *AlertService) Get(ctx context.Context, id int64) (*model.Alert, error) {
	return s.store.GetByID(ctx, id)
}

func (s *AlertService) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.Alert, error) {
	return s.store.ListByIncubator(ctx, incubatorID)
}

func (s *AlertService) ListByLevel(ctx context.Context, level string) ([]*model.Alert, error) {
	return s.store.ListByLevel(ctx, level)
}

func (s *AlertService) Acknowledge(ctx context.Context, id int64) error {
	return s.store.Acknowledge(ctx, id)
}

func (s *AlertService) CountActive(ctx context.Context, incubatorID int64) (int, error) {
	return s.store.CountActive(ctx, incubatorID)
}
