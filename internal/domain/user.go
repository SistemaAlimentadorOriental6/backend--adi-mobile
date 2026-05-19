package domain

// Usuario representa la entidad de usuario en el sistema
type Usuario struct {
	Cedula string `json:"cedula"`
	Nombre string `json:"nombre"`
	Cargo  string `json:"cargo"`
	Email  string `json:"email"`
}

// UserRepository define el contrato para el acceso a datos de usuario
type UserRepository interface {
	GetByCedula(cedula string, email string) (*Usuario, error)
	Upsert(user *Usuario) error
	GetActiveEmployees() (map[string]Usuario, error)
}
