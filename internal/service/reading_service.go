package service

import (
	"context"
	"sync"
	"time"

	"github.com/jb843051627/hatchery/internal/cache"
	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type ReadingService struct {
	store *store.ReadingStore
	cache *cache.ReadingCache
}

func NewReadingService(s *store.ReadingStore, rc *cache.ReadingCache) *ReadingService {
	return &ReadingService{store: s, cache: rc}
}

func (s *ReadingService) Record(ctx context.Context, incubatorID int64, sensorType string, value float64) (int64, error) {
	now := time.Now().UTC()
	id, err := s.store.Create(ctx, incubatorID, sensorType, value, now)
	if err != nil {
		return 0, err
	}
	s.cache.Update(&model.SensorReading{
		ID:          id,
		IncubatorID: incubatorID,
		SensorType:  sensorType,
		Value:       value,
		RecordedAt:  now,
	})
	return id, nil
}

func (s *ReadingService) Get(ctx context.Context, id int64) (*model.SensorReading, error) {
	return s.store.GetByID(ctx, id)
}

func (s *ReadingService) ListByIncubator(ctx context.Context, incubatorID int64, limit int) ([]*model.SensorReading, error) {
	return s.store.ListByIncubator(ctx, incubatorID, limit)
}

func (s *ReadingService) ListByTimeRange(ctx context.Context, incubatorID int64, from, to time.Time) ([]*model.SensorReading, error) {
	return s.store.ListByTimeRange(ctx, incubatorID, from, to)
}

func (s *ReadingService) BatchIngest(ctx context.Context, readings []*model.SensorReading) (int, error) {
	count := 0
	var mu sync.Mutex
	for _, r := range readings {
		if ctx.Err() != nil {
			break
		}
		id, err := s.store.Create(ctx, r.IncubatorID, r.SensorType, r.Value, r.RecordedAt)
		if err != nil {
			continue
		}
		s.cache.Update(&model.SensorReading{
			ID:          id,
			IncubatorID: r.IncubatorID,
			SensorType:  r.SensorType,
			Value:       r.Value,
			RecordedAt:  r.RecordedAt,
		})
		mu.Lock()
		count++
		mu.Unlock()
	}
	return count, nil
}

func (s *ReadingService) LatestByIncubator(ctx context.Context, incubatorID int64, sensorType string) (*model.SensorReading, error) {
	return s.store.LatestByIncubator(ctx, incubatorID, sensorType)
}

func (s *ReadingService) CountByIncubator(ctx context.Context, incubatorID int64) (int, error) {
	return s.store.CountByIncubator(ctx, incubatorID)
}
