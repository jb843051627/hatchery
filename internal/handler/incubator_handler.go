package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/hatchery/internal/service"
)

type IncubatorHandler struct {
	svc *service.IncubatorService
}

func NewIncubatorHandler(s *service.IncubatorService) *IncubatorHandler {
	return &IncubatorHandler{svc: s}
}

func (h *IncubatorHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := writeIDParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	inc, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

func (h *IncubatorHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	var list interface{}
	var err error
	if status != "" {
		list, err = h.svc.ListByStatus(r.Context(), status)
	} else {
		list, err = h.svc.ListAll(r.Context())
	}
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *IncubatorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Capacity int    `json:"capacity"`
		Status   string `json:"status"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	id, err := h.svc.Create(r.Context(), req.Name, req.Location, req.Capacity, req.Status)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *IncubatorHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
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

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
