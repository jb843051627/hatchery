package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/hatchery/internal/cache"
	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/store"
)

type ReportService struct {
	st *store.Store
	rc *cache.ReadingCache
}

func NewReportService(st *store.Store, rc *cache.ReadingCache) *ReportService {
	return &ReportService{st: st, rc: rc}
}

type IncubatorSummary struct {
	Incubator    *model.Incubator     `json:"incubator"`
	ActiveBatch  *model.Batch         `json:"active_batch"`
	LatestReadings map[string]*model.SensorReading `json:"latest_readings"`
	ActiveAlerts int                  `json:"active_alerts"`
}

func (s *ReportService) IncubatorSummary(ctx context.Context, incubatorID int64) (*IncubatorSummary, error) {
	inc, err := s.st.Incubators.GetByID(ctx, incubatorID)
	if err != nil {
		return nil, fmt.Errorf("get incubator: %w", err)
	}
	if inc == nil {
		return nil, model.ErrIncubatorNotFound
	}
	batch, _ := s.st.Batches.GetActiveBatch(ctx, incubatorID)
	alertCount, _ := s.st.Alerts.CountActive(ctx, incubatorID)
	readings := s.rc.AllForIncubator(incubatorID)
	return &IncubatorSummary{
		Incubator:      inc,
		ActiveBatch:    batch,
		LatestReadings: readings,
		ActiveAlerts:   alertCount,
	}, nil
}

type BatchReport struct {
	Batch       *model.Batch         `json:"batch"`
	Trays       []*model.EggTray     `json:"trays"`
	Schedules   []*model.Schedule    `json:"schedules"`
	HatchRecords []*model.HatchRecord `json:"hatch_records"`
	TotalHatched int                 `json:"total_hatched"`
	HealthyRate  float64             `json:"healthy_rate"`
}

func (s *ReportService) BatchReport(ctx context.Context, batchID int64) (*BatchReport, error) {
	batch, err := s.st.Batches.GetByID(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("get batch: %w", err)
	}
	if batch == nil {
		return nil, model.ErrBatchNotFound
	}
	trays, _ := s.st.EggTrays.ListByBatch(ctx, batchID)
	schedules, _ := s.st.Schedules.ListByBatch(ctx, batchID)
	records, _ := s.st.HatchRecords.ListByBatch(ctx, batchID)

	out := &BatchReport{
		Batch:       batch,
		Trays:       trays,
		Schedules:   schedules,
		HatchRecords: records,
	}
	if records != nil {
		total, healthy, _, _, _ := s.st.HatchRecords.SummaryByBatch(ctx, batchID)
		out.TotalHatched = total
		if total > 0 {
			out.HealthyRate = float64(healthy) / float64(total) * 100
		}
	}
	return out, nil
}

type ExportRow struct {
	IncubatorID int64   `json:"incubator_id"`
	SensorType  string  `json:"sensor_type"`
	Value       float64 `json:"value"`
	RecordedAt  string  `json:"recorded_at"`
}

func (s *ReportService) ExportReadings(ctx context.Context, incubatorID int64, from, to time.Time) ([]*ExportRow, error) {
	readings, err := s.st.Readings.ListByTimeRange(ctx, incubatorID, from, to)
	if err != nil {
		return nil, fmt.Errorf("list readings: %w", err)
	}
	out := make([]*ExportRow, len(readings))
	for i, r := range readings {
		out[i] = &ExportRow{
			IncubatorID: r.IncubatorID,
			SensorType:  r.SensorType,
			Value:       r.Value,
			RecordedAt:  r.RecordedAt.Format("2006-01-02 15:04:05"),
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].RecordedAt < out[j].RecordedAt
	})
	loc, _ := time.LoadLocation("Asia/Shanghai")
	if loc != nil {
		for i, row := range out {
			if t, err := time.ParseInLocation("2006-01-02 15:04:05", row.RecordedAt, loc); err == nil {
				out[i].RecordedAt = t.Format("2006-01-02 15:04:05")
			}
		}
	}
	return out, nil
}

type DashboardStats struct {
	TotalIncubators  int            `json:"total_incubators"`
	ActiveBatches    int            `json:"active_batches"`
	ActiveAlerts     int            `json:"active_alerts"`
	PendingMaintenance int          `json:"pending_maintenance"`
	RecentHatchRecords []*model.HatchRecord `json:"recent_hatch_records"`
}

func (s *ReportService) Dashboard(ctx context.Context) (*DashboardStats, error) {
	incubatorCount, _ := s.st.Incubators.Count(ctx)
	activeBatches, _ := s.st.Batches.CountByStatus(ctx, model.BatchStatusIncubating)
	pendingMaintenance, _ := s.st.Maintenance.ListPending(ctx)
	records, _ := s.st.HatchRecords.ListByDateRange(ctx, time.Now().UTC().AddDate(0, -1, 0), time.Now().UTC())

	return &DashboardStats{
		TotalIncubators:    incubatorCount,
		ActiveBatches:      activeBatches,
		PendingMaintenance: len(pendingMaintenance),
		RecentHatchRecords: records,
	}, nil
}
