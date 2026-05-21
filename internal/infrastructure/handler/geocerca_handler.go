package handler

import (
	"encoding/json"
	"net/http"
	"turnovia-backend/internal/usecase"
)

// GeocercaHandler expone las rutas HTTP correspondientes a las geocercas.
type GeocercaHandler struct {
	casoUso usecase.GeocercaUseCase
}

// NewGeocercaHandler inicializa el handler de geocercas.
func NewGeocercaHandler(uc usecase.GeocercaUseCase) *GeocercaHandler {
	return &GeocercaHandler{casoUso: uc}
}

// Listar procesa la petición GET para retornar la lista detallada de geocercas.
func (h *GeocercaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Método no permitido, se requiere GET",
		})
		return
	}

	geocercas, err := h.casoUso.ObtenerGeocercas()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error interno al recuperar las geocercas: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    geocercas,
	})
}
