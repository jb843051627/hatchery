package handler

import (
	"net/http"
	"strconv"

	"github.com/jb843051627/hatchery/internal/service"
)

type EggTrayHandler struct {
	svc *service.EggTrayService
}

func NewEggTrayHandler(s *service.EggTrayService) *EggTrayHandler {
	return &EggTrayHandler{svc: s}
}

func (h *EggTrayHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BatchID    int64   `json:"batch_id"`
		TrayNumber int     `json:"tray_number"`
		EggCount   int     `json:"egg_count"`
		Weight     float64 `json:"weight"`
		SourceFarm string  `json:"source_farm"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	id, err := h.svc.Create(r.Context(), req.BatchID, req.TrayNumber, req.EggCount, req.Weight, req.SourceFarm)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *EggTrayHandler) ListByBatch(w http.ResponseWriter, r *http.Request) {
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
