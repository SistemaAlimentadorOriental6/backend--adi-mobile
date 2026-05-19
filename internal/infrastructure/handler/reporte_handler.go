package handler

import (
	"encoding/json"
	"net/http"
	"turnovia-backend/internal/usecase"
)

type ReporteHandler struct {
	useCase usecase.ReporteUseCase
}

func NewReporteHandler(u usecase.ReporteUseCase) *ReporteHandler {
	return &ReporteHandler{useCase: u}
}

type reporteRequest struct {
	Cedula          string  `json:"cedula"`
	Nombre          string  `json:"nombre"`
	EntradaSalida   string  `json:"entradasalida"`
	Lugar           string  `json:"lugar"`
	Latitud         float64 `json:"latitud"`
	Longitud        float64 `json:"longitud"`
	Precision       string  `json:"precision"`
	Dispositivo     string  `json:"dispositivo"`
	TipoRed         string  `json:"tipo_red"`
	MetodoUbicacion string  `json:"metodo_ubicacion"`
	ZonaHoraria     string  `json:"zona_horaria"`
}

func (h *ReporteHandler) Registrar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req reporteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if req.EntradaSalida == "" {
		req.EntradaSalida = "entrada"
	}

	err := h.useCase.CrearReporte(
		req.Cedula,
		req.Nombre,
		req.EntradaSalida,
		req.Lugar,
		req.Latitud,
		req.Longitud,
		req.Dispositivo,
		req.TipoRed,
		req.MetodoUbicacion,
		req.ZonaHoraria,
	)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al guardar el reporte: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Reporte registrado exitosamente",
	})
}

func (h *ReporteHandler) Listar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	reportes, err := h.useCase.ObtenerReportes(50)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al obtener reportes: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    reportes,
	})
}

func (h *ReporteHandler) ObtenerRegistroHoy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	cedula := r.URL.Query().Get("cedula")
	if cedula == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Cédula requerida",
		})
		return
	}

	registro, err := h.useCase.ObtenerRegistroHoy(cedula)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al obtener registro: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    registro,
	})
}
