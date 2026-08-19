package service

import (
	"github.com/jb843051627/hatchery/internal/cache"
	"github.com/jb843051627/hatchery/internal/store"
)

type Service struct {
	Incubators   *IncubatorService
	Batches      *BatchService
	EggTrays     *EggTrayService
	Readings     *ReadingService
	Schedules    *ScheduleService
	HatchRecords *HatchRecordService
	Alerts       *AlertService
	Maintenance  *MaintenanceService
	Report       *ReportService
}

func NewService(st *store.Store, rc *cache.ReadingCache) *Service {
	return &Service{
		Incubators:   NewIncubatorService(st.Incubators),
		Batches:      NewBatchService(st.Batches, st.Incubators),
		EggTrays:     NewEggTrayService(st.EggTrays),
		Readings:     NewReadingService(st.Readings, rc),
		Schedules:    NewScheduleService(st.Schedules),
		HatchRecords: NewHatchRecordService(st.HatchRecords),
		Alerts:       NewAlertService(st.Alerts),
		Maintenance:  NewMaintenanceService(st.Maintenance, st.DB()),
		Report:       NewReportService(st, rc),
	}
}
