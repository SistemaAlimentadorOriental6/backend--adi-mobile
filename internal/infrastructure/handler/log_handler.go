package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"turnovia-backend/internal/usecase"
)

type LogHandler struct {
	useCase usecase.LogUseCase
}

func NewLogHandler(u usecase.LogUseCase) *LogHandler {
	return &LogHandler{useCase: u}
}

type logBatchRequest struct {
	Logs []usecase.LogItem `json:"logs"`
}

func (h *LogHandler) RegistrarBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req logBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if len(req.Logs) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"count":   0,
		})
		return
	}

	if err := h.useCase.RegistrarBatch(req.Logs); err != nil {
		log.Printf("❌ Error guardando logs: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al guardar logs",
		})
		return
	}

	log.Printf("📋 Logs recibidos: %d entradas", len(req.Logs))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(req.Logs),
	})
}
