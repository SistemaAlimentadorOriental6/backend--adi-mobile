package domain

import "time"

type TrackingUbicacion struct {
	Id        uint64    `json:"id"`
	Cedula    string    `json:"cedula"`
	Lugar     string    `json:"lugar"`
	Latitud   float64   `json:"latitud"`
	Longitud  float64   `json:"longitud"`
	Timestamp time.Time `json:"timestamp"`
	Estado    string    `json:"estado"`
	Validado  bool      `json:"validado"`
}

type TrackingRepository interface {
	Save(t *TrackingUbicacion) error
	SaveBatch(items []*TrackingUbicacion) error
	HasActiveSession(cedula string) (bool, error)
}