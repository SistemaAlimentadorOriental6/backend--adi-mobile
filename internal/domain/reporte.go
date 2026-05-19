package domain

import "time"

// Reporte representa un registro de actividad con geolocalización
type Reporte struct {
	ID              int       `json:"id"`
	Cedula          string    `json:"cedula"`
	Nombre          string    `json:"nombre"`
	EntradaSalida   string    `json:"entradasalida"` // 'entrada' o 'salida'
	Lugar           string    `json:"lugar"`
	Latitud         float64   `json:"latitud"`
	Longitud        float64   `json:"longitud"`
	Tiempo          time.Time `json:"tiempo"`
	CreatedAt       string    `json:"created_at"`
	Dispositivo     string    `json:"dispositivo"`
	TipoRed         string    `json:"tipo_red"`
	MetodoUbicacion string    `json:"metodo_ubicacion"`
	ZonaHoraria     string    `json:"zona_horaria"`
}

// RegistroHoy representa el estado de registro del día actual
type RegistroHoy struct {
	HayEntrada bool   `json:"hay_entrada"`
	HaySalida  bool   `json:"hay_salida"`
	Lugar      string `json:"lugar"`
}

// ReporteRepository define el contrato para guardar y listar los reportes
type ReporteRepository interface {
	Save(reporte *Reporte) error
	ListarReportes(limit int) ([]*Reporte, error)
	GetRegistroHoy(cedula string) (*RegistroHoy, error)
}
