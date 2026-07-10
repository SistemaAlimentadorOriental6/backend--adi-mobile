package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"turnovia-backend/internal/usecase"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(u usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: u}
}

type loginRequest struct {
	Cedula string `json:"cedula"`
	Email  string `json:"email"`
}

type loginResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	User    *interface{} `json:"user,omitempty"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// CORS header - Aunque ya se maneja en el middleware, se mantiene por seguridad
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var cedula string
	log.Printf("📥 Petición recibida: %s %s", r.Method, r.URL.String())

	var email string
	if r.Method == "POST" {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("❌ Error decodificando body: %v", err)
			cedula = r.URL.Query().Get("cedula")
			email = r.URL.Query().Get("email")
		} else {
			cedula = req.Cedula
			email = req.Email
		}
	} else if r.Method == "GET" {
		cedula = r.URL.Query().Get("cedula")
		email = r.URL.Query().Get("email")
		log.Printf("ℹ️ Datos GET: cedula='%s', email='%s' (Query: %s)", cedula, email, r.URL.RawQuery)
	} else {
		log.Printf("❌ Método no permitido: %s", r.Method)
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if cedula == "" || email == "" {
		log.Printf("⚠️ Intento de login sin datos completos. URL: %s", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "La cédula y el email son requeridos",
		})
		return
	}

	user, err := h.authUseCase.Login(cedula, email)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "¡Bienvenido de nuevo!",
		"user":    user,
	})
}

type registerBiometricsRequest struct {
	Cedula string `json:"cedula"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

func (h *AuthHandler) RegisterBiometrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Método no permitido"})
		return
	}

	var req registerBiometricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Cuerpo inválido"})
		return
	}

	if req.Cedula == "" || req.Email == "" || req.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "La cédula, el email y el token son requeridos"})
		return
	}

	err := h.authUseCase.RegisterBiometrics(req.Cedula, req.Email, req.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token biométrico registrado correctamente",
	})
}

type loginBiometricsRequest struct {
	Token string `json:"token"`
}

func (h *AuthHandler) LoginBiometrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Método no permitido"})
		return
	}

	var req loginBiometricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Cuerpo inválido"})
		return
	}

	if req.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "El token es requerido"})
		return
	}

	user, err := h.authUseCase.LoginWithBiometrics(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "¡Bienvenido de nuevo!",
		"user":    user,
	})
}

type removeBiometricsRequest struct {
	Cedula string `json:"cedula"`
}

func (h *AuthHandler) RemoveBiometrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Método no permitido"})
		return
	}

	var req removeBiometricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Cuerpo inválido"})
		return
	}

	if req.Cedula == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "La cédula es requerida"})
		return
	}

	err := h.authUseCase.RemoveBiometrics(req.Cedula)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token biométrico removido correctamente",
	})
}

