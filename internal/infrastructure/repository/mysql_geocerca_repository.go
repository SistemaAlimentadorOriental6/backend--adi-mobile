package repository

import (
	"database/sql"
	"strconv"
	"turnovia-backend/internal/domain"
)

type mysqlGeocercaRepository struct {
	baseDatos *sql.DB
}

// NewMysqlGeocercaRepository crea una instancia del repositorio de geocercas para MySQL.
func NewMysqlGeocercaRepository(db *sql.DB) domain.GeocercaRepository {
	return &mysqlGeocercaRepository{baseDatos: db}
}

// ObtenerTodas recupera las geocercas de MySQL ordenadas por nombre y orden_punto.
func (r *mysqlGeocercaRepository) ObtenerTodas() ([]*domain.Geocerca, error) {
	// Consulta los puntos ordenados de cada ubicación geográfica
	consulta := `SELECT COALESCE(nombre, ''), COALESCE(orden_punto, 0), COALESCE(latitude, ''), COALESCE(longitude, '') 
		FROM geocercas 
		ORDER BY nombre, orden_punto`

	filas, err := r.baseDatos.Query(consulta)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	geocercasMapa := make(map[string][]domain.Punto)
	var nombres []string

	for filas.Next() {
		var nombre string
		var ordenPunto int
		var latitudTexto string
		var longitudTexto string

		err := filas.Scan(&nombre, &ordenPunto, &latitudTexto, &longitudTexto)
		if err != nil {
			return nil, err
		}

		// Conversión de las coordenadas guardadas como texto a float64
		latitud, errLat := strconv.ParseFloat(latitudTexto, 64)
		longitud, errLng := strconv.ParseFloat(longitudTexto, 64)
		if errLat != nil || errLng != nil {
			continue
		}

		punto := domain.Punto{
			Latitud:  latitud,
			Longitud: longitud,
		}

		if _, existe := geocercasMapa[nombre]; !existe {
			nombres = append(nombres, nombre)
		}
		geocercasMapa[nombre] = append(geocercasMapa[nombre], punto)
	}

	var resultado []*domain.Geocerca
	for _, nombre := range nombres {
		resultado = append(resultado, &domain.Geocerca{
			Nombre: nombre,
			Puntos: geocercasMapa[nombre],
		})
	}

	return resultado, nil
}
