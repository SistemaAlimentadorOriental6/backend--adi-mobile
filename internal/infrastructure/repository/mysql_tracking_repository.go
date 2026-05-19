package repository

import (
	"database/sql"
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