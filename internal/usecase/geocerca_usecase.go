package usecase

import "turnovia-backend/internal/domain"

// GeocercaUseCase expone las operaciones de negocio permitidas para las geocercas.
type GeocercaUseCase interface {
	ObtenerGeocercas() ([]*domain.Geocerca, error)
}

type geocercaUseCase struct {
	repositorio domain.GeocercaRepository
}

// NewGeocercaUseCase inicializa el caso de uso de geocercas.
func NewGeocercaUseCase(repo domain.GeocercaRepository) GeocercaUseCase {
	return &geocercaUseCase{repositorio: repo}
}

func (uc *geocercaUseCase) ObtenerGeocercas() ([]*domain.Geocerca, error) {
	return uc.repositorio.ObtenerTodas()
}
