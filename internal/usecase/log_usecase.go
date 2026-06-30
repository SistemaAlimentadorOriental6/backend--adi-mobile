package usecase

import (
	"time"
	"turnovia-backend/internal/domain"
)

// LogItem es el DTO que recibe el handler desde el móvil
type LogItem struct {
	Cedula      string `json:"cedula"`
	Nivel       string `json:"nivel"`
	Tag         string `json:"tag"`
	Mensaje     string `json:"mensaje"`
	Extra       string `json:"extra"`
	Dispositivo string `json:"dispositivo"`
	VersionApp  string `json:"version_app"`
	Timestamp   int64  `json:"timestamp"` // milisegundos desde epoch
}

type LogUseCase interface {
	RegistrarBatch(items []LogItem) error
}

type logUseCase struct {
	repo domain.LogRepository
}

func NewLogUseCase(repo domain.LogRepository) LogUseCase {
	return &logUseCase{repo: repo}
}

func (u *logUseCase) RegistrarBatch(items []LogItem) error {
	if len(items) == 0 {
		return nil
	}

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		loc = time.UTC
	}

	logs := make([]*domain.AppLog, 0, len(items))
	for _, item := range items {
		nivel := domain.NivelLog(item.Nivel)
		if nivel == "" {
			nivel = domain.NivelInfo
		}

		logs = append(logs, &domain.AppLog{
			Cedula:      item.Cedula,
			Nivel:       nivel,
			Tag:         item.Tag,
			Mensaje:     item.Mensaje,
			Extra:       item.Extra,
			Dispositivo: item.Dispositivo,
			VersionApp:  item.VersionApp,
			Timestamp:   time.UnixMilli(item.Timestamp).In(loc),
		})
	}

	return u.repo.GuardarBatch(logs)
}
