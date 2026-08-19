package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/service"
)

type HatchHandler struct {
	svc *service.HatchRecordService
}

func NewHatchHandler(s *service.HatchRecordService) *HatchHandler {
	return &HatchHandler{svc: s}
}

func (h *HatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BatchID      int64  `json:"batch_id"`
		HatchedCount int    `json:"hatched_count"`
		HealthyCount int    `json:"healthy_count"`
		WeakCount    int    `json:"weak_count"`
		DeadCount    int    `json:"dead_count"`
		HatchDate    string `json:"hatch_date"`
		GradedBy     string `json:"graded_by"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	hatchDate, err := time.Parse(time.RFC3339, req.HatchDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid hatch_date"})
		return
	}
	rec := &model.HatchRecord{
		BatchID:      req.BatchID,
		HatchedCount: req.HatchedCount,
		HealthyCount: req.HealthyCount,
		WeakCount:    req.WeakCount,
		DeadCount:    req.DeadCount,
		HatchDate:    hatchDate,
		GradedBy:     req.GradedBy,
	}
	id, err := h.svc.Create(r.Context(), rec)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *HatchHandler) ListByBatch(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "batch_id")
	batchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid batch_id"})
		return
	}
	list, err := h.svc.ListByBatch(r.Context(), batchID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
