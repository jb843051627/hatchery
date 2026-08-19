package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/hatchery/internal/service"
)

type BatchHandler struct {
	svc *service.BatchService
}

func NewBatchHandler(s *service.BatchService) *BatchHandler {
	return &BatchHandler{svc: s}
}

func (h *BatchHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	b, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *BatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IncubatorID   int64  `json:"incubator_id"`
		EggCount      int    `json:"egg_count"`
		Species       string `json:"species"`
		StartDate     string `json:"start_date"`
		ExpectedHatch string `json:"expected_hatch"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid start_date"})
		return
	}
	expectedHatch, err := time.Parse(time.RFC3339, req.ExpectedHatch)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid expected_hatch"})
		return
	}
	id, err := h.svc.Create(r.Context(), req.IncubatorID, req.EggCount, req.Species, startDate, expectedHatch)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *BatchHandler) ListByIncubator(w http.ResponseWriter, r *http.Request) {
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

func (h *BatchHandler) StartIncubation(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.svc.StartIncubation(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *BatchHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := h.svc.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
