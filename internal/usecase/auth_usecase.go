package usecase

import (
	"log"
	"turnovia-backend/internal/domain"
)

type AuthUseCase interface {
	Login(cedula string, email string) (*domain.Usuario, error)
	RegisterBiometrics(cedula string, email string, token string) error
	LoginWithBiometrics(token string) (*domain.Usuario, error)
	RemoveBiometrics(cedula string) error
}

type authUseCase struct {
	mysqlRepo     domain.UserRepository
	sqlServerRepo domain.UserRepository
}

// NewAuthUseCase crea una nueva instancia de la lógica de negocio de autenticación
func NewAuthUseCase(mysqlRepo domain.UserRepository, sqlServerRepo domain.UserRepository) AuthUseCase {
	return &authUseCase{
		mysqlRepo:     mysqlRepo,
		sqlServerRepo: sqlServerRepo,
	}
}

func (u *authUseCase) Login(cedula string, email string) (*domain.Usuario, error) {
	log.Printf("🔍 Intentando login para cédula: %s, email: %s", cedula, email)

	// Obtener Nombre, Cargo, Email y validar acceso desde SQL Server directamente
	userDetails, err := u.sqlServerRepo.GetByCedula(cedula, email)
	if err != nil {
		log.Printf("❌ Error en SQL Server para %s: %v", cedula, err)
		return nil, err
	}
	// 2. Sincronizar el usuario en MySQL para evitar errores de llave foránea al guardar reportes
	err = u.mysqlRepo.Upsert(userDetails)
	if err != nil {
		log.Printf("⚠️ No se pudo sincronizar el usuario en MySQL: %v", err)
		// No exponemos el error al login porque ya está validado por SQL Server
	} else {
		log.Printf("✅ Usuario %s sincronizado en MySQL", cedula)
	}

	return userDetails, nil
}

func (u *authUseCase) RegisterBiometrics(cedula string, email string, token string) error {
	log.Printf("🔒 Registrando token biométrico para la cédula: %s", cedula)
	return u.mysqlRepo.SaveBiometricToken(cedula, email, token)
}

func (u *authUseCase) LoginWithBiometrics(token string) (*domain.Usuario, error) {
	log.Printf("🔍 Intentando login con token biométrico")
	cachedUser, err := u.mysqlRepo.GetByBiometricToken(token)
	if err != nil {
		log.Printf("❌ Error buscando token biométrico en MySQL: %v", err)
		return nil, err
	}

	// Validar el usuario contra SQL Server (origen de verdad) usando su cédula y correo guardados
	log.Printf("🔍 Validando usuario biométrico en SQL Server para cédula: %s, email: %s", cachedUser.Cedula, cachedUser.Email)
	userDetails, err := u.sqlServerRepo.GetByCedula(cachedUser.Cedula, cachedUser.Email)
	if err != nil {
		log.Printf("❌ Validación biométrica falló en SQL Server para %s: %v", cachedUser.Cedula, err)
		return nil, err
	}

	// Sincronizar de nuevo los datos en MySQL
	_ = u.mysqlRepo.Upsert(userDetails)

	log.Printf("✅ Login biométrico exitoso y verificado en SQL Server para: %s (%s)", userDetails.Nombre, userDetails.Cedula)
	return userDetails, nil
}

func (u *authUseCase) RemoveBiometrics(cedula string) error {
	log.Printf("🔓 Desvinculando token biométrico para la cédula: %s", cedula)
	return u.mysqlRepo.RemoveBiometricToken(cedula)
}
