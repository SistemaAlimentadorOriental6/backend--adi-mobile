package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"turnovia-backend/internal/domain"
)

type ProxyHandler struct {
	repo     domain.ProxyRepository
	userRepo domain.UserRepository
}

func NewProxyHandler(repo domain.ProxyRepository, userRepo domain.UserRepository) *ProxyHandler {
	return &ProxyHandler{repo: repo, userRepo: userRepo}
}

func (h *ProxyHandler) GetEmpleados(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener empleados de la API externa
	apiData, err := h.repo.GetEmpleados()
	if err != nil {
		h.sendError(w, err)
		return
	}

	// 2. Obtener empleados activos de SQL Server
	activeMap, err := h.userRepo.GetActiveEmployees()
	if err != nil {
		h.sendError(w, err)
		return
	}

	// 3. Filtrar y enriquecer
	var filtered []domain.Empleado
	log.Printf("📊 Cruzando %d empleados de API con %d activos en SQL Server", len(apiData.Data), len(activeMap))

	for _, emp := range apiData.Data {
		cleanCedula := strings.TrimSpace(emp.Cedula)
		if sqlUser, ok := activeMap[cleanCedula]; ok {
			// Enriquecer con datos de SQL Server
			emp.Nombres = sqlUser.Nombre // Usamos el nombre de SQL Server
			emp.Cargo = sqlUser.Cargo    // Usamos el cargo de SQL Server
			filtered = append(filtered, emp)
		}
	}

	log.Printf("✅ %d operadores finales después de filtrar activos", len(filtered))

	// 4. Responder con los datos filtrados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.EmpleadoApiResponse{
		Status: true,
		Count:  len(filtered),
		Data:   filtered,
	})
}

func (h *ProxyHandler) sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  false,
		"message": err.Error(),
	})
}
