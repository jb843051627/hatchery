package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/service"
)

type MaintenanceHandler struct {
	svc *service.MaintenanceService
}

func NewMaintenanceHandler(s *service.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{svc: s}
}

func (h *MaintenanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IncubatorID  int64  `json:"incubator_id"`
		Description  string `json:"description"`
		ScheduledFor string `json:"scheduled_for"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	scheduledFor, err := time.Parse(time.RFC3339, req.ScheduledFor)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid scheduled_for"})
		return
	}
	id, err := h.svc.Create(r.Context(), req.IncubatorID, req.Description, scheduledFor)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *MaintenanceHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListPending(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MaintenanceHandler) Complete(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.svc.Complete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MaintenanceHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	var reqs []struct {
		IncubatorID  int64  `json:"incubator_id"`
		Description  string `json:"description"`
		ScheduledFor string `json:"scheduled_for"`
	}
	if err := decodeJSON(r, &reqs); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	tasks := make([]*model.MaintenanceTask, 0, len(reqs))
	for _, req := range reqs {
		scheduledFor, err := time.Parse(time.RFC3339, req.ScheduledFor)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid scheduled_for"})
			return
		}
		tasks = append(tasks, &model.MaintenanceTask{
			IncubatorID:  req.IncubatorID,
			Description:  req.Description,
			ScheduledFor: scheduledFor,
		})
	}
	if err := h.svc.BatchCreate(r.Context(), tasks); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int{"count": len(tasks)})
}
