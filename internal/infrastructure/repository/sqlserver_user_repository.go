package repository

import (
	"database/sql"
	"errors"
	"turnovia-backend/internal/domain"
)

type sqlServerUserRepository struct {
	db *sql.DB
}

// NewSqlServerUserRepository crea una instancia del repositorio de usuario para SQL Server
func NewSqlServerUserRepository(db *sql.DB) domain.UserRepository {
	return &sqlServerUserRepository{db: db}
}

func (r *sqlServerUserRepository) GetByCedula(cedula string, email string) (*domain.Usuario, error) {
	query := `
		SELECT TOP 1 f_nombre_empl, f_desc_cargo, f_email_contacto
		FROM SE_W0550
		WHERE f_nit_empl = @p1 AND LTRIM(RTRIM(f_email_contacto)) = LTRIM(RTRIM(@p2))
		ORDER BY f_parametro DESC, f_ndc DESC
	`

	user := &domain.Usuario{Cedula: cedula}
	var cargo, emailContacto string
	err := r.db.QueryRow(query, cedula, email).Scan(&user.Nombre, &cargo, &emailContacto)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado o email incorrecto")
		}
		return nil, err
	}

	cargosAutorizados := map[string]bool{
		"OPERADOR MASTER":         true,
		"OPERADOR":                true,
		"AUXILIAR DE INTEGRACION": true,
		"REGULADOR VIA":           true,
		"AUXILIAR LOGISTICO":       true,
		"ANALISTA DE MEJORA CONTINUA Y AUTOMATIZACION": true,
		"PROFESIONAL DE PLANEACION Y PROGRAMACION":     true,
	}

	if !cargosAutorizados[cargo] {
		return nil, errors.New("el cargo '" + cargo + "' no tiene permisos para acceder a esta aplicación")
	}

	user.Cargo = cargo
	user.Email = emailContacto
	return user, nil
}

func (r *sqlServerUserRepository) GetActiveEmployees() (map[string]domain.Usuario, error) {
	// Query con TRIM para eliminar espacios en blanco de SQL Server
	query := `
		WITH ActiveEmpl AS (
			SELECT 
				LTRIM(RTRIM(f_nit_empl)) as cedula, 
				LTRIM(RTRIM(f_nombre_empl)) as nombre, 
				LTRIM(RTRIM(f_desc_cargo)) as cargo,
				ROW_NUMBER() OVER(PARTITION BY f_nit_empl ORDER BY f_parametro DESC, f_ndc DESC) as rn
			FROM SE_W0550 
			WHERE (f_fecha_retiro IS NULL OR f_fecha_retiro = '' OR f_fecha_retiro = '1900-01-01')
			  AND LTRIM(RTRIM(f_desc_cargo)) IN (
				'OPERADOR', 
				'OPERADOR DE ALISTAMIENTO', 
				'OPERADOR DE DIAGNOSTICO', 
				'OPERADOR DUAL', 
				'OPERADOR EN FORMACION', 
				'OPERADOR MASTER'
				'AUXILIAR DE INTEGRACION',
				'REGULADOR VIA',
				'AUXILIAR DE FLOTA',
				'AUXILIAR LOGISTICO',
			  )
		)
		SELECT cedula, nombre, cargo 
		FROM ActiveEmpl 
		WHERE rn = 1
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activeMap := make(map[string]domain.Usuario)
	for rows.Next() {
		var u domain.Usuario
		var cargoNull sql.NullString // Manejar posibles nulos

		if err := rows.Scan(&u.Cedula, &u.Nombre, &cargoNull); err != nil {
			continue
		}

		if cargoNull.Valid {
			u.Cargo = cargoNull.String
		} else {
			u.Cargo = "OPERADOR" // Valor por defecto si es nulo
		}

		// Guardar en el mapa con la cédula limpia
		activeMap[u.Cedula] = u
	}

	return activeMap, nil
}

func (r *sqlServerUserRepository) Upsert(user *domain.Usuario) error {
	return nil
}

func (r *sqlServerUserRepository) SaveBiometricToken(cedula string, token string) error {
	return nil
}

func (r *sqlServerUserRepository) GetByBiometricToken(token string) (*domain.Usuario, error) {
	return nil, errors.New("autenticación biométrica no soportada en SQL Server")
}
