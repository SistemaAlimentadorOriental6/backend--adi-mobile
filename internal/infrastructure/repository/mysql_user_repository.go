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

func (r *mysqlUserRepository) SaveBiometricToken(cedula string, email string, token string) error {
	query := `
		INSERT INTO validacion_biometrico (cedula, email, biometric_token) 
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE email = VALUES(email), biometric_token = VALUES(biometric_token)
	`
	_, err := r.db.Exec(query, cedula, email, token)
	return err
}

func (r *mysqlUserRepository) RemoveBiometricToken(cedula string) error {
	query := "DELETE FROM validacion_biometrico WHERE cedula = ?"
	_, err := r.db.Exec(query, cedula)
	return err
}

func (r *mysqlUserRepository) GetByBiometricToken(token string) (*domain.Usuario, error) {
	query := `
		SELECT a.cedula, a.nombre, COALESCE(a.cargo, ''), vb.email, vb.biometric_token 
		FROM validacion_biometrico vb
		INNER JOIN auxiliares a ON vb.cedula = a.cedula
		WHERE vb.biometric_token = ?
	`
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
	// Asegurar tabla auxiliares
	createAuxiliaresQuery := `
		CREATE TABLE IF NOT EXISTS auxiliares (
			cedula VARCHAR(50) PRIMARY KEY,
			nombre VARCHAR(150) NOT NULL,
			cargo VARCHAR(100) DEFAULT NULL
		);
	`
	if _, err := r.db.Exec(createAuxiliaresQuery); err != nil {
		return err
	}

	// Asegurar tabla validacion_biometrico
	createBiometricoQuery := `
		CREATE TABLE IF NOT EXISTS validacion_biometrico (
			cedula VARCHAR(50) PRIMARY KEY,
			email VARCHAR(150) NOT NULL,
			biometric_token VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uq_biometric_token (biometric_token)
		);
	`
	if _, err := r.db.Exec(createBiometricoQuery); err != nil {
		return err
	}

	// Por si la tabla auxiliares antigua no tenía cargo
	r.db.Exec("ALTER TABLE auxiliares ADD COLUMN cargo VARCHAR(100) DEFAULT NULL")
	return nil
}
