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
	Estado       string  `json:"estado"`
}

type TrackingUseCase interface {
	GuardarUbicacion(cedula, lugar string, lat, lng float64, estado string, validado bool) error
	GuardarUbicacionesBatch(items []TrackingBatchItem) error
}

type trackingUseCase struct {
	repo         domain.TrackingRepository
	geocercaRepo domain.GeocercaRepository
}

func NewTrackingUseCase(repo domain.TrackingRepository, geocercaRepo domain.GeocercaRepository) TrackingUseCase {
	return &trackingUseCase{repo: repo, geocercaRepo: geocercaRepo}
}

func (u *trackingUseCase) GuardarUbicacion(cedula, lugar string, lat, lng float64, estado string, validado bool) error {
	hasActive, err := u.repo.HasActiveSession(cedula)
	if err != nil {
		return err
	}
	if !hasActive {
		return ErrNoActiveSession
	}

	// Validar usando geocercas de base de datos si no fue validada previamente o como validación definitiva
	if lugar != "" {
		geocercas, errG := u.geocercaRepo.ObtenerTodas()
		if errG == nil {
			for _, g := range geocercas {
				if g.Nombre == lugar {
					validado = domain.PointInPolygon(lat, lng, g.Puntos)
					break
				}
			}
		}
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

	// Cargar todas las geocercas una única vez para validar el lote completo en memoria
	mapaGeocercas := make(map[string][]domain.Punto)
	geocercas, errG := u.geocercaRepo.ObtenerTodas()
	if errG == nil {
		for _, g := range geocercas {
			mapaGeocercas[g.Nombre] = g.Puntos
		}
	}

	// Extraer todas las marcas de tiempo para consultar y filtrar los duplicados en la base de datos
	var tsList []time.Time
	for _, it := range items {
		loc, _ := time.LoadLocation("America/Bogota")
		if loc == nil {
			loc = time.UTC
		}
		ts := time.UnixMilli(it.Timestamp).In(loc)
		tsList = append(tsList, ts)
	}

	// Obtener marcas de tiempo que ya existen en MySQL para evitar inserciones duplicadas (idempotencia)
	existingMap, err := u.repo.GetExistingTimestamps(items[0].Cedula, tsList)
	if err != nil {
		return err
	}

	var toSave []*domain.TrackingUbicacion
	for _, it := range items {
		// Convertir timestamp de JS (milisegundos) a time.Time en Bogotá
		loc, _ := time.LoadLocation("America/Bogota")
		if loc == nil {
			loc = time.UTC
		}
		ts := time.UnixMilli(it.Timestamp).In(loc)

		// Filtrar si ya existe en la base de datos MySQL (por reintento de red del cliente)
		key := ts.Format("2006-01-02 15:04:05.000")
		if existingMap[key] {
			continue
		}

		// Calcular si el punto está dentro de la geocerca para este lugar
		validadoItem := it.Validado
		if it.Lugar != "" {
			if puntos, existe := mapaGeocercas[it.Lugar]; existe {
				validadoItem = domain.PointInPolygon(it.Latitud, it.Longitud, puntos)
			}
		}

		estadoItem := it.Estado
		if estadoItem == "" {
			estadoItem = "ok"
		}

		toSave = append(toSave, &domain.TrackingUbicacion{
			Cedula:    it.Cedula,
			Lugar:     it.Lugar,
			Latitud:   it.Latitud,
			Longitud:  it.Longitud,
			Timestamp: ts,
			Estado:    estadoItem,
			Validado:  validadoItem,
		})
	}

	// Si todos los elementos del batch ya existían, no es necesario llamar a SaveBatch
	if len(toSave) == 0 {
		return nil
	}

	return u.repo.SaveBatch(toSave)
}
