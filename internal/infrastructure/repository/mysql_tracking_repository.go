package repository

import (
	"database/sql"
	"time"
	"turnovia-backend/internal/domain"
)

type mysqlTrackingRepository struct {
	db *sql.DB
}

func NewMysqlTrackingRepository(db *sql.DB) domain.TrackingRepository {
	return &mysqlTrackingRepository{db: db}
}

func (r *mysqlTrackingRepository) Save(t *domain.TrackingUbicacion) error {
	query := `INSERT INTO tracking_ubicacion (cedula, lugar, latitud, longitud, timestamp, estado, validado) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		t.Cedula,
		t.Lugar,
		t.Latitud,
		t.Longitud,
		t.Timestamp,
		t.Estado,
		t.Validado,
	)

	return err
}

func (r *mysqlTrackingRepository) SaveBatch(items []*domain.TrackingUbicacion) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO tracking_ubicacion (cedula, lugar, latitud, longitud, timestamp, estado, validado) VALUES (?, ?, ?, ?, ?, ?, ?)`

	for _, t := range items {
		_, err := tx.Exec(query,
			t.Cedula,
			t.Lugar,
			t.Latitud,
			t.Longitud,
			t.Timestamp,
			t.Estado,
			t.Validado,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *mysqlTrackingRepository) HasActiveSession(cedula string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM registros 
		WHERE cedula = ? 
		AND entradasalida = 'entrada' 
		AND DATE(tiempo) = CURDATE()
		AND NOT EXISTS (
			SELECT 1 FROM registros r2 
			WHERE r2.cedula = registros.cedula 
			AND r2.entradasalida = 'salida' 
			AND DATE(r2.tiempo) = CURDATE()
		)
	`
	row := r.db.QueryRow(query, cedula)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetExistingTimestamps consulta a la base de datos MySQL por marcas de tiempo de tracking que ya existen para este usuario.
func (r *mysqlTrackingRepository) GetExistingTimestamps(cedula string, timestamps []time.Time) (map[string]bool, error) {
	if len(timestamps) == 0 {
		return make(map[string]bool), nil
	}

	query := `SELECT timestamp FROM tracking_ubicacion WHERE cedula = ? AND timestamp IN (`
	args := []interface{}{cedula}
	for i, ts := range timestamps {
		if i > 0 {
			query += ","
		}
		query += "?"
		args = append(args, ts)
	}
	query += ")"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]bool)
	for rows.Next() {
		var t time.Time
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		// Formatear llave para comparación segura en milisegundos
		key := t.Format("2006-01-02 15:04:05.000")
		existing[key] = true
	}

	return existing, nil
}