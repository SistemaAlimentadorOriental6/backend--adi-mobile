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
	query := "SELECT cedula, nombre FROM auxiliares WHERE cedula = ?"
	row := r.db.QueryRow(query, cedula)

	user := &domain.Usuario{}
	err := row.Scan(&user.Cedula, &user.Nombre)
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
		INSERT INTO auxiliares (cedula, nombre) 
		VALUES (?, ?) 
		ON DUPLICATE KEY UPDATE nombre = VALUES(nombre)
	`
	_, err := r.db.Exec(query, user.Cedula, user.Nombre)
	return err
}

func (r *mysqlUserRepository) GetActiveEmployees() (map[string]domain.Usuario, error) {
	return make(map[string]domain.Usuario), nil
}
