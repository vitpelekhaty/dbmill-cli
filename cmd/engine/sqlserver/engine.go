package sqlserver

import (
	"database/sql"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

// Engine реализация функциональности утилиты dbmill-cli для MS SQL Server
type Engine struct {
	db     *sql.DB
	logger *log.Logger
}

// NewEngine возвращает экземпляр Engine
func NewEngine(connection string) (*Engine, error) {
	return nil, nil
}

func (self *Engine) SetLogger(logger *log.Logger) {
	self.logger = logger
}
