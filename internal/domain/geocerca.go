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

// PointInPolygon determina si un punto (latitud, longitud) está dentro del polígono de la geocerca.
func PointInPolygon(latitud, longitud float64, poligono []Punto) bool {
	dentro := false
	n := len(poligono)
	if n < 3 {
		return false
	}
	j := n - 1
	for i := 0; i < n; i++ {
		xi := poligono[i].Latitud
		yi := poligono[i].Longitud
		xj := poligono[j].Latitud
		yj := poligono[j].Longitud

		interseccion := ((yi > longitud) != (yj > longitud)) &&
			(latitud < (xj-xi)*(longitud-yi)/(yj-yi)+xi)
		if interseccion {
			dentro = !dentro
		}
		j = i
	}
	return dentro
}

