package usecase

import (
	"time"
	"turnovia-backend/internal/domain"
)

type ReporteUseCase interface {
	CrearReporte(cedula, nombre, entradasalida, lugar string, lat, lng float64, dispositivo, red, metodo, zona string) error
	ObtenerReportes(limit int) ([]*domain.Reporte, error)
	ObtenerRegistroHoy(cedula string) (*domain.RegistroHoy, error)
}

type reporteUseCase struct {
	repo domain.ReporteRepository
}

func NewReporteUseCase(repo domain.ReporteRepository) ReporteUseCase {
	return &reporteUseCase{repo: repo}
}

func (u *reporteUseCase) CrearReporte(cedula, nombre, entradasalida, lugar string, lat, lng float64, dispositivo, red, metodo, zona string) error {
	reporte := &domain.Reporte{
		Cedula:          cedula,
		Nombre:          nombre,
		EntradaSalida:   entradasalida,
		Lugar:           lugar,
		Latitud:         lat,
		Longitud:        lng,
		Tiempo:          time.Now(),
		Dispositivo:     dispositivo,
		TipoRed:         red,
		MetodoUbicacion: metodo,
		ZonaHoraria:     zona,
	}
	return u.repo.Save(reporte)
}

func (u *reporteUseCase) ObtenerReportes(limit int) ([]*domain.Reporte, error) {
	return u.repo.ListarReportes(limit)
}

func (u *reporteUseCase) ObtenerRegistroHoy(cedula string) (*domain.RegistroHoy, error) {
	return u.repo.GetRegistroHoy(cedula)
}
