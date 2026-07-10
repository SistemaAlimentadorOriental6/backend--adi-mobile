package repository

import (
	"database/sql"
	"errors"
	"turnovia-backend/internal/domain"
)

type mysqlUserRepository struct {
	db *sql.DB
}

// NewMysqlUserRepository crea una instancia del repositorio de usuario para MySQL
func NewMysqlUserRepository(db *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetByCedula(cedula string, email string) (*domain.Usuario, error) {
	query := "SELECT cedula, nombre, COALESCE(cargo, ''), COALESCE(email, ''), COALESCE(biometric_token, '') FROM auxiliares WHERE cedula = ?"
	row := r.db.QueryRow(query, cedula)

	user := &domain.Usuario{}
	err := row.Scan(&user.Cedula, &user.Nombre, &user.Cargo, &user.Email, &user.BiometricToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return user, nil
}

func (r *mysqlUserRepository) Upsert(user *domain.Usuario) error {
	query := `
		INSERT INTO auxiliares (cedula, nombre, cargo, email) 
		VALUES (?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE nombre = VALUES(nombre), cargo = VALUES(cargo), email = VALUES(email)
	`
	_, err := r.db.Exec(query, user.Cedula, user.Nombre, user.Cargo, user.Email)
	return err
}

func (r *mysqlUserRepository) GetActiveEmployees() (map[string]domain.Usuario, error) {
	return make(map[string]domain.Usuario), nil
}

func (r *mysqlUserRepository) SaveBiometricToken(cedula string, token string) error {
	query := "UPDATE auxiliares SET biometric_token = ? WHERE cedula = ?"
	_, err := r.db.Exec(query, token, cedula)
	return err
}

func (r *mysqlUserRepository) GetByBiometricToken(token string) (*domain.Usuario, error) {
	query := "SELECT cedula, nombre, COALESCE(cargo, ''), COALESCE(email, ''), COALESCE(biometric_token, '') FROM auxiliares WHERE biometric_token = ?"
	row := r.db.QueryRow(query, token)

	user := &domain.Usuario{}
	err := row.Scan(&user.Cedula, &user.Nombre, &user.Cargo, &user.Email, &user.BiometricToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("token biométrico inválido o no registrado")
		}
		return nil, err
	}

	return user, nil
}

func (r *mysqlUserRepository) EnsureSchema() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS auxiliares (
			cedula VARCHAR(50) PRIMARY KEY,
			nombre VARCHAR(150) NOT NULL,
			cargo VARCHAR(100) DEFAULT NULL,
			email VARCHAR(150) DEFAULT NULL,
			biometric_token VARCHAR(255) DEFAULT NULL
		);
	`
	if _, err := r.db.Exec(createTableQuery); err != nil {
		return err
	}

	r.db.Exec("ALTER TABLE auxiliares ADD COLUMN cargo VARCHAR(100) DEFAULT NULL")
	r.db.Exec("ALTER TABLE auxiliares ADD COLUMN email VARCHAR(150) DEFAULT NULL")
	r.db.Exec("ALTER TABLE auxiliares ADD COLUMN biometric_token VARCHAR(255) DEFAULT NULL")
	return nil
}
