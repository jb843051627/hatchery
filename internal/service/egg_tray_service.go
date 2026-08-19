package service

import (
	"context"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type EggTrayService struct {
	store *store.EggTrayStore
}

func NewEggTrayService(s *store.EggTrayStore) *EggTrayService {
	return &EggTrayService{store: s}
}

func (s *EggTrayService) Create(ctx context.Context, batchID int64, trayNumber int, eggCount int, weight float64, sourceFarm string) (int64, error) {
	return s.store.Create(ctx, batchID, trayNumber, eggCount, weight, sourceFarm)
}

func (s *EggTrayService) Get(ctx context.Context, id int64) (*model.EggTray, error) {
	t, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, model.ErrEggTrayNotFound
	}
	return t, nil
}

func (s *EggTrayService) ListByBatch(ctx context.Context, batchID int64) ([]*model.EggTray, error) {
	return s.store.ListByBatch(ctx, batchID)
}

func (s *EggTrayService) UpdateEggCount(ctx context.Context, id int64, count int) error {
	return s.store.UpdateEggCount(ctx, id, count)
}

func (s *EggTrayService) DeleteByBatch(ctx context.Context, batchID int64) error {
	return s.store.DeleteByBatch(ctx, batchID)
}

func (s *EggTrayService) CountByBatch(ctx context.Context, batchID int64) (int, error) {
	return s.store.CountByBatch(ctx, batchID)
}

func (s *EggTrayService) SumEggsByBatch(ctx context.Context, batchID int64) (int, error) {
	return s.store.SumEggsByBatch(ctx, batchID)
}
