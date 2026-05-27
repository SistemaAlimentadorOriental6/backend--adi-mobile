package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"turnovia-backend/internal/usecase"
)

type TrackingHandler struct {
	useCase usecase.TrackingUseCase
}

func NewTrackingHandler(u usecase.TrackingUseCase) *TrackingHandler {
	return &TrackingHandler{useCase: u}
}

type trackingRequest struct {
	Cedula   string  `json:"cedula"`
	Lugar    string  `json:"lugar"`
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
	Validado bool    `json:"validado"`
	Estado   string  `json:"estado"`
}

type trackingBatchItemRequest struct {
	Latitud      float64 `json:"latitud"`
	Longitud     float64 `json:"longitud"`
	Accuracy     float64 `json:"accuracy"`
	Timestamp    int64   `json:"timestamp"`
	IsStationary bool    `json:"is_stationary"`
	Lugar        string  `json:"lugar"`
	Cedula       string  `json:"cedula"`
	Estado       string  `json:"estado"`
}

type trackingBatchRequest struct {
	Locations []trackingBatchItemRequest `json:"locations"`
}

func (h *TrackingHandler) GuardarUbicacionBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req trackingBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if len(req.Locations) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "El batch está vacío",
		})
		return
	}

	log.Printf("📍 Batch recibido: %d ubicaciones", len(req.Locations))

	var batchItems []usecase.TrackingBatchItem
	for _, loc := range req.Locations {
		estado := loc.Estado
		if estado == "" {
			estado = "ok"
		}
		batchItems = append(batchItems, usecase.TrackingBatchItem{
			Cedula:       loc.Cedula,
			Lugar:        loc.Lugar,
			Latitud:      loc.Latitud,
			Longitud:     loc.Longitud,
			Accuracy:     loc.Accuracy,
			Timestamp:    loc.Timestamp,
			IsStationary: loc.IsStationary,
			Validado:     false,
			Estado:       estado,
		})
	}

	err := h.useCase.GuardarUbicacionesBatch(batchItems)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, usecase.ErrNoActiveSession) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Sesión no iniciada. Registra una entrada primero.",
			})
			return
		}
		log.Printf("❌ Error guardando batch: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al guardar ubicaciones",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(req.Locations),
	})
}

func (h *TrackingHandler) GuardarUbicacion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req trackingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	log.Printf("📍 Tracking recibido: cedula=%s lugar=%s lat=%f lng=%f validado=%v", req.Cedula, req.Lugar, req.Latitud, req.Longitud, req.Validado)

	estado := req.Estado
	if estado == "" {
		estado = "ok"
	}
	err := h.useCase.GuardarUbicacion(req.Cedula, req.Lugar, req.Latitud, req.Longitud, estado, req.Validado)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, usecase.ErrNoActiveSession) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Sesión no iniciada. Registra una entrada primero.",
			})
			return
		}
		log.Printf("❌ Error guardando tracking: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al guardar ubicación",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
