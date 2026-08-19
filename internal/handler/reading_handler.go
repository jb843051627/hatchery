package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
	"github.com/jb843051627/hatchery/internal/service"
)

type ReadingHandler struct {
	svc *service.ReadingService
}

func NewReadingHandler(s *service.ReadingService) *ReadingHandler {
	return &ReadingHandler{svc: s}
}

func (h *ReadingHandler) Record(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IncubatorID int64   `json:"incubator_id"`
		SensorType  string  `json:"sensor_type"`
		Value       float64 `json:"value"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	id, err := h.svc.Record(r.Context(), req.IncubatorID, req.SensorType, req.Value)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *ReadingHandler) BatchIngest(w http.ResponseWriter, r *http.Request) {
	var readings []*model.SensorReading
	if err := decodeJSON(r, &readings); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	count, err := h.svc.BatchIngest(r.Context(), readings)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *ReadingHandler) Export(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	incubatorIDStr := q.Get("incubator_id")
	incubatorID, err := strconv.ParseInt(incubatorIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid incubator_id"})
		return
	}
	from, err := time.Parse(time.RFC3339, q.Get("from"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid from"})
		return
	}
	to, err := time.Parse(time.RFC3339, q.Get("to"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid to"})
		return
	}
	readings, err := h.svc.ListByTimeRange(r.Context(), incubatorID, from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, readings)
}
