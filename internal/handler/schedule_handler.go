package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/service"
)

type ScheduleHandler struct {
	svc *service.ScheduleService
}

func NewScheduleHandler(s *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: s}
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IncubatorID    int64   `json:"incubator_id"`
		BatchID        int64   `json:"batch_id"`
		Phase          string  `json:"phase"`
		TargetTemp     float64 `json:"target_temp"`
		TargetHumidity float64 `json:"target_humidity"`
		StartTime      string  `json:"start_time"`
		EndTime        string  `json:"end_time"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid start_time"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid end_time"})
		return
	}
	sc := &model.Schedule{
		IncubatorID:    req.IncubatorID,
		BatchID:        req.BatchID,
		Phase:          req.Phase,
		TargetTemp:     req.TargetTemp,
		TargetHumidity: req.TargetHumidity,
		StartTime:      startTime,
		EndTime:        endTime,
	}
	id, err := h.svc.Create(r.Context(), sc)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *ScheduleHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	sc, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sc)
}

func (h *ScheduleHandler) Activate(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.svc.Activate(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *ScheduleHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.svc.Cancel(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
