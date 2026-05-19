package usecase

import (
	"errors"
	"time"
	"turnovia-backend/internal/domain"
)

var ErrNoActiveSession = errors.New("no existe una sesión activa para este usuario")

func getBogotaTime() time.Time {
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return time.Now().UTC().Add(-5 * time.Hour)
	}
	return time.Now().In(loc)
}

type TrackingBatchItem struct {
	Cedula       string  `json:"cedula"`
	Lugar        string  `json:"lugar"`
	Latitud      float64 `json:"latitud"`
	Longitud     float64 `json:"longitud"`
	Accuracy     float64 `json:"accuracy"`
	Timestamp    int64   `json:"timestamp"`
	IsStationary bool    `json:"is_stationary"`
	Validado     bool    `json:"validado"`
}

type TrackingUseCase interface {
	GuardarUbicacion(cedula, lugar string, lat, lng float64, estado string, validado bool) error
	GuardarUbicacionesBatch(items []TrackingBatchItem) error
}

type trackingUseCase struct {
	repo domain.TrackingRepository
}

func NewTrackingUseCase(repo domain.TrackingRepository) TrackingUseCase {
	return &trackingUseCase{repo: repo}
}

func (u *trackingUseCase) GuardarUbicacion(cedula, lugar string, lat, lng float64, estado string, validado bool) error {
	hasActive, err := u.repo.HasActiveSession(cedula)
	if err != nil {
		return err
	}
	if !hasActive {
		return ErrNoActiveSession
	}

	tracking := &domain.TrackingUbicacion{
		Cedula:    cedula,
		Lugar:     lugar,
		Latitud:   lat,
		Longitud:  lng,
		Timestamp: getBogotaTime(),
		Estado:    estado,
		Validado:  validado,
	}
	return u.repo.Save(tracking)
}

func (u *trackingUseCase) GuardarUbicacionesBatch(items []TrackingBatchItem) error {
	if len(items) == 0 {
		return nil
	}

	// Verificar sesión activa una sola vez (todo el batch viene del mismo usuario)
	hasActive, err := u.repo.HasActiveSession(items[0].Cedula)
	if err != nil {
		return err
	}
	if !hasActive {
		return ErrNoActiveSession
	}

	var toSave []*domain.TrackingUbicacion
	for _, it := range items {
		// Convertir timestamp de JS (milisegundos) a time.Time en Bogotá
		loc, _ := time.LoadLocation("America/Bogota")
		if loc == nil {
			loc = time.UTC
		}
		ts := time.UnixMilli(it.Timestamp).In(loc)

		toSave = append(toSave, &domain.TrackingUbicacion{
			Cedula:    it.Cedula,
			Lugar:     it.Lugar,
			Latitud:   it.Latitud,
			Longitud:  it.Longitud,
			Timestamp: ts,
			Estado:    "ok",
			Validado:  it.Validado,
		})
	}

	return u.repo.SaveBatch(toSave)
}