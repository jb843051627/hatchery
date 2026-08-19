package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jb843051627/hatchery/internal/service"
)

type AlertHandler struct {
	svc *service.AlertService
}

func NewAlertHandler(s *service.AlertService) *AlertHandler {
	return &AlertHandler{svc: s}
}

func (h *AlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IncubatorID int64  `json:"incubator_id"`
		BatchID     int64  `json:"batch_id"`
		Level       string `json:"level"`
		Message     string `json:"message"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	id, err := h.svc.Create(r.Context(), req.IncubatorID, req.BatchID, req.Level, req.Message)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *AlertHandler) ListByIncubator(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "incubator_id")
	incubatorID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid incubator_id"})
		return
	}
	list, err := h.svc.ListByIncubator(r.Context(), incubatorID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *AlertHandler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.svc.Acknowledge(r.Context(), id); err != nil {
		writeError(w, fmt.Errorf("acknowledge alert: %w", err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
