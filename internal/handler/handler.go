package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/service"
)

type Handler struct {
	Incubators   *IncubatorHandler
	Batches      *BatchHandler
	EggTrays     *EggTrayHandler
	Readings     *ReadingHandler
	Schedules    *ScheduleHandler
	HatchRecords *HatchHandler
	Alerts       *AlertHandler
	Maintenance  *MaintenanceHandler
	Report       *ReportHandler
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		Incubators:   NewIncubatorHandler(svc.Incubators),
		Batches:      NewBatchHandler(svc.Batches),
		EggTrays:     NewEggTrayHandler(svc.EggTrays),
		Readings:     NewReadingHandler(svc.Readings),
		Schedules:    NewScheduleHandler(svc.Schedules),
		HatchRecords: NewHatchHandler(svc.HatchRecords),
		Alerts:       NewAlertHandler(svc.Alerts),
		Maintenance:  NewMaintenanceHandler(svc.Maintenance),
		Report:       NewReportHandler(svc.Report),
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/incubators", h.Incubators.List)
	mux.HandleFunc("GET /api/incubators/{id}", h.Incubators.Get)
	mux.HandleFunc("POST /api/incubators", h.Incubators.Create)
	mux.HandleFunc("PUT /api/incubators/{id}/status", h.Incubators.UpdateStatus)
	mux.HandleFunc("POST /api/batches", h.Batches.Create)
	mux.HandleFunc("GET /api/batches/{id}", h.Batches.Get)
	mux.HandleFunc("GET /api/incubators/{incubator_id}/batches", h.Batches.ListByIncubator)
	mux.HandleFunc("PUT /api/batches/{id}/start", h.Batches.StartIncubation)
	mux.HandleFunc("PUT /api/batches/{id}/status", h.Batches.UpdateStatus)
	mux.HandleFunc("POST /api/egg-trays", h.EggTrays.Create)
	mux.HandleFunc("GET /api/batches/{batch_id}/egg-trays", h.EggTrays.ListByBatch)
	mux.HandleFunc("POST /api/readings", h.Readings.Record)
	mux.HandleFunc("POST /api/readings/batch", h.Readings.BatchIngest)
	mux.HandleFunc("GET /api/readings/export", h.Readings.Export)
	mux.HandleFunc("POST /api/schedules", h.Schedules.Create)
	mux.HandleFunc("GET /api/schedules/{id}", h.Schedules.Get)
	mux.HandleFunc("PUT /api/schedules/{id}/activate", h.Schedules.Activate)
	mux.HandleFunc("PUT /api/schedules/{id}/cancel", h.Schedules.Cancel)
	mux.HandleFunc("POST /api/hatch-records", h.HatchRecords.Create)
	mux.HandleFunc("GET /api/batches/{batch_id}/hatch-records", h.HatchRecords.ListByBatch)
	mux.HandleFunc("POST /api/alerts", h.Alerts.Create)
	mux.HandleFunc("GET /api/incubators/{incubator_id}/alerts", h.Alerts.ListByIncubator)
	mux.HandleFunc("PUT /api/alerts/{id}/ack", h.Alerts.Acknowledge)
	mux.HandleFunc("POST /api/maintenance", h.Maintenance.Create)
	mux.HandleFunc("GET /api/maintenance/pending", h.Maintenance.ListPending)
	mux.HandleFunc("PUT /api/maintenance/{id}/complete", h.Maintenance.Complete)
	mux.HandleFunc("POST /api/maintenance/batch", h.Maintenance.BatchCreate)
	mux.HandleFunc("GET /api/reports/dashboard", h.Report.Dashboard)
	mux.HandleFunc("GET /api/reports/incubator/{id}", h.Report.IncubatorSummary)
	mux.HandleFunc("GET /api/reports/batch/{id}", h.Report.BatchReport)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return mux
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var ve *model.ValidationError
	if errors.As(err, &ve) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": ve.Error(), "field": ve.Field})
		return
	}
	sentinels := []error{
		model.ErrIncubatorNotFound, model.ErrBatchNotFound, model.ErrEggTrayNotFound,
		model.ErrReadingNotFound, model.ErrScheduleNotFound, model.ErrHatchRecordNotFound,
		model.ErrAlertNotFound, model.ErrMaintenanceNotFound,
	}
	for _, sentinel := range sentinels {
		if errors.Is(err, sentinel) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func writeCreated(w http.ResponseWriter, id int64) {
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func writeIDParam(r *http.Request, name string) string {
	return r.PathValue(name)
}
