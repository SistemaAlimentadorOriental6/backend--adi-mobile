package repository

import (
	"database/sql"
	"turnovia-backend/internal/domain"
)

// MysqlLogRepository exportado para permitir acceso a EnsureSchema desde main.go
type MysqlLogRepository struct {
	db *sql.DB
}

func NewMysqlLogRepository(db *sql.DB) domain.LogRepository {
	return &MysqlLogRepository{db: db}
}

// GuardarBatch inserta un lote de logs en MySQL dentro de una sola transacción.
// No usa INSERT IGNORE porque los logs casi nunca se duplican y no tienen clave natural única clara.
func (r *MysqlLogRepository) GuardarBatch(logs []*domain.AppLog) error {
	if len(logs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO logs 
			(cedula, nivel, tag, mensaje, extra, dispositivo, version_app, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, l := range logs {
		_, err := stmt.Exec(
			l.Cedula,
			string(l.Nivel),
			l.Tag,
			l.Mensaje,
			l.Extra,
			l.Dispositivo,
			l.VersionApp,
			l.Timestamp,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// EnsureSchema crea la tabla logs si no existe. Llamar al arrancar el servidor.
func (r *MysqlLogRepository) EnsureSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS logs (
			id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			cedula       VARCHAR(30)  NOT NULL DEFAULT '',
			nivel        VARCHAR(10)  NOT NULL,
			tag          VARCHAR(60)  NOT NULL,
			mensaje      TEXT         NOT NULL,
			extra        TEXT         NULL,
			dispositivo  VARCHAR(120) NULL,
			version_app  VARCHAR(20)  NULL,
			timestamp    DATETIME(3)  NOT NULL,
			creado_en    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			INDEX idx_cedula_ts (cedula, timestamp),
			INDEX idx_nivel     (nivel),
			INDEX idx_tag       (tag)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`
	_, err := r.db.Exec(query)
	return err
}
