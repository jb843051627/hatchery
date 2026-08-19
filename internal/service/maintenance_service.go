package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type MaintenanceService struct {
	store *store.MaintenanceStore
}

func NewMaintenanceService(s *store.MaintenanceStore, _ interface{}) *MaintenanceService {
	return &MaintenanceService{store: s}
}

func (s *MaintenanceService) Create(ctx context.Context, incubatorID int64, description string, scheduledFor time.Time) (int64, error) {
	return s.store.Create(ctx, incubatorID, description, scheduledFor)
}

func (s *MaintenanceService) Get(ctx context.Context, id int64) (*model.MaintenanceTask, error) {
	return s.store.GetByID(ctx, id)
}

func (s *MaintenanceService) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.MaintenanceTask, error) {
	return s.store.ListByIncubator(ctx, incubatorID)
}

func (s *MaintenanceService) ListPending(ctx context.Context) ([]*model.MaintenanceTask, error) {
	return s.store.ListPending(ctx)
}

func (s *MaintenanceService) Complete(ctx context.Context, id int64) error {
	t, err := s.store.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get maintenance task: %w", err)
	}
	if t == nil {
		return model.ErrMaintenanceNotFound
	}
	if t.Status == model.MaintenanceStatusCompleted {
		return &model.ValidationError{Field: "status", Message: "task already completed"}
	}
	return s.store.Complete(ctx, id, time.Now().UTC())
}

func (s *MaintenanceService) BatchCreate(ctx context.Context, tasks []*model.MaintenanceTask) error {
	if len(tasks) == 0 {
		return &model.ValidationError{Field: "tasks", Message: "tasks is empty"}
	}
	return s.store.BatchCreate(ctx, tasks)
}
