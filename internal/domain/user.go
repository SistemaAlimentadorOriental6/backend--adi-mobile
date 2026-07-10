package domain

// Usuario representa la entidad de usuario en el sistema
type Usuario struct {
	Cedula         string `json:"cedula"`
	Nombre         string `json:"nombre"`
	Cargo          string `json:"cargo"`
	Email          string `json:"email"`
	BiometricToken string `json:"biometric_token,omitempty"`
}

// UserRepository define el contrato para el acceso a datos de usuario
type UserRepository interface {
	GetByCedula(cedula string, email string) (*Usuario, error)
	Upsert(user *Usuario) error
	GetActiveEmployees() (map[string]Usuario, error)
	SaveBiometricToken(cedula string, email string, token string) error
	GetByBiometricToken(token string) (*Usuario, error)
	RemoveBiometricToken(cedula string) error
}
