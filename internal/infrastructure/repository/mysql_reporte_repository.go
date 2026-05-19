package repository

import (
	"database/sql"
	"turnovia-backend/internal/domain"
)

type mysqlReporteRepository struct {
	db *sql.DB
}

// NewMysqlReporteRepository crea una instancia del repositorio de reporte para MySQL
func NewMysqlReporteRepository(db *sql.DB) domain.ReporteRepository {
	return &mysqlReporteRepository{db: db}
}

func (r *mysqlReporteRepository) Save(reporte *domain.Reporte) error {
	query := `INSERT INTO registros 
		(cedula, nombre, entradasalida, lugar, latitud, longitud, dispositivo, tipo_red, metodo_ubicacion, zona_horaria, tiempo) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		reporte.Cedula,
		reporte.Nombre,
		reporte.EntradaSalida,
		reporte.Lugar,
		reporte.Latitud,
		reporte.Longitud,
		reporte.Dispositivo,
		reporte.TipoRed,
		reporte.MetodoUbicacion,
		reporte.ZonaHoraria,
		reporte.Tiempo)

	return err
}

func (r *mysqlReporteRepository) ListarReportes(limit int) ([]*domain.Reporte, error) {
	query := `SELECT id, cedula, nombre, entradasalida, lugar, latitud, longitud, dispositivo, tipo_red, metodo_ubicacion, zona_horaria, tiempo, created_at 
		FROM registros 
		WHERE DATE(created_at) = CURDATE()
		ORDER BY created_at DESC 
		LIMIT ?`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportes []*domain.Reporte
	for rows.Next() {
		var r domain.Reporte
		err := rows.Scan(
			&r.ID,
			&r.Cedula,
			&r.Nombre,
			&r.EntradaSalida,
			&r.Lugar,
			&r.Latitud,
			&r.Longitud,
			&r.Dispositivo,
			&r.TipoRed,
			&r.MetodoUbicacion,
			&r.ZonaHoraria,
			&r.Tiempo,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reportes = append(reportes, &r)
	}

	return reportes, nil
}

func (r *mysqlReporteRepository) GetRegistroHoy(cedula string) (*domain.RegistroHoy, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN entradasalida = 'entrada' THEN 1 END), 0) > 0 AS hay_entrada,
			COALESCE(SUM(CASE WHEN entradasalida = 'salida' THEN 1 END), 0) > 0 AS hay_salida,
			COALESCE(MAX(CASE WHEN entradasalida = 'entrada' THEN lugar END), '') AS lugar
		FROM registros 
		WHERE cedula = ? AND DATE(tiempo) = CURDATE()
	`
	row := r.db.QueryRow(query, cedula)

	var result domain.RegistroHoy
	err := row.Scan(&result.HayEntrada, &result.HaySalida, &result.Lugar)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
