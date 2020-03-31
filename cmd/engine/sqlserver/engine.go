package sqlserver

import (
	"database/sql"
)

// Engine реализация функциональности утилиты dbmill-cli для MS SQL Server
type Engine struct {
	db     *sql.DB
	logger interface{}
}

// NewEngine возвращает экземпляр Engine
func NewEngine(connection string) (*Engine, error) {
	return nil, nil
}

func (self *Engine) SetLogger(logger interface{}) {
	self.logger = logger
}
