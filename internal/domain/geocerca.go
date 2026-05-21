package domain

import "time"

// Punto representa una coordenada geográfica individual de una geocerca.
type Punto struct {
	Latitud  float64 `json:"latitude"`
	Longitud float64 `json:"longitude"`
}

// Geocerca agrupa un conjunto de puntos ordenados bajo un mismo nombre identificador.
type Geocerca struct {
	Nombre string  `json:"nombre"`
	Puntos []Punto `json:"coords"`
}

// GeocercaRegistro modela la estructura directa de una fila en la tabla de base de datos.
type GeocercaRegistro struct {
	ID         int       `json:"id"`
	Nombre     string    `json:"nombre"`
	OrdenPunto int       `json:"orden_punto"`
	Latitud    string    `json:"latitude"`
	Longitud   string    `json:"longitude"`
	CreadoEn   time.Time `json:"created_at"`
}

// GeocercaRepository define el contrato para acceder a los datos de geocercas en el sistema.
type GeocercaRepository interface {
	ObtenerTodas() ([]*Geocerca, error)
}
