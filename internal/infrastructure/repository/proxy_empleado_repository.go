package repository

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"turnovia-backend/internal/domain"
)

type proxyEmpleadoRepository struct {
	apiURL string
	apiKey string
	client *http.Client
}

func NewProxyEmpleadoRepository() domain.ProxyRepository {
	// Configuramos transporte para ignorar errores de certificado SSL si los hay
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &proxyEmpleadoRepository{
		apiURL: "https://adbordo.valliu.co:8443/api/employee",
		apiKey: "7caa01b9ef1c2b902efadf6211add0ca41bdae9f",
		client: &http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
		},
	}
}

func (r *proxyEmpleadoRepository) GetEmpleados() (*domain.EmpleadoApiResponse, error) {
	log.Printf("🌐 Consultando API externa: %s", r.apiURL)

	req, err := http.NewRequest("GET", r.apiURL, nil)
	if err != nil {
		log.Printf("❌ Error creando petición: %v", err)
		return nil, err
	}

	// Go canonicaliza los headers con Set() → "X-Api-Key". La API exige "x-api-key" en minúsculas.
	req.Header["x-api-key"] = []string{r.apiKey}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "PostmanRuntime/7.43.0")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := r.client.Do(req)
	if err != nil {
		log.Printf("❌ Error en petición HTTP: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("⚠️ API externa respondió con status: %d — Cuerpo: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("error api externa: %d — %s", resp.StatusCode, string(bodyBytes))
	}

	var result domain.EmpleadoApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("❌ Error decodificando respuesta JSON: %v", err)
		return nil, err
	}

	log.Printf("✅ %d empleados obtenidos exitosamente", len(result.Data))
	return &result, nil
}
