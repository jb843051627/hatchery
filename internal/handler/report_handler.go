package handler

import (
	"net/http"
	"strconv"

	"github.com/jb843051627/hatchery/internal/service"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: s}
}

func (h *ReportHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.Dashboard(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *ReportHandler) IncubatorSummary(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	summary, err := h.svc.IncubatorSummary(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *ReportHandler) BatchReport(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	report, err := h.svc.BatchReport(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}
