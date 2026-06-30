package domain

import "time"

// NivelLog define la severidad de un evento registrado
type NivelLog string

const (
	NivelDebug NivelLog = "DEBUG"
	NivelInfo  NivelLog = "INFO"
	NivelWarn  NivelLog = "WARN"
	NivelError NivelLog = "ERROR"
	NivelFatal NivelLog = "FATAL"
)

// AppLog representa un evento de diagnóstico generado por la app móvil
type AppLog struct {
	Id          uint64    `json:"id"`
	Cedula      string    `json:"cedula"`
	Nivel       NivelLog  `json:"nivel"`
	Tag         string    `json:"tag"`
	Mensaje     string    `json:"mensaje"`
	Extra       string    `json:"extra"`
	Dispositivo string    `json:"dispositivo"`
	VersionApp  string    `json:"version_app"`
	Timestamp   time.Time `json:"timestamp"`
}

type LogRepository interface {
	GuardarBatch(logs []*AppLog) error
}
