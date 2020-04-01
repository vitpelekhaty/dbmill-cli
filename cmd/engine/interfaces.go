package engine

import (
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

// IDatabase интерфейс базы данных
type IDatabase interface {
	SetLogger(logger *log.Logger)
	ScriptsFolder(path string, includeData, decrypt bool) error
}
