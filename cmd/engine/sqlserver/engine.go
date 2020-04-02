package sqlserver

import (
	"database/sql"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/dir"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

// Engine реализация функциональности утилиты dbmill-cli для MS SQL Server
type Engine struct {
	db      *sql.DB
	logger  log.ILogger
	include filter.IFilter
	exclude filter.IFilter
	output  dir.IStructure
}

// NewEngine возвращает экземпляр Engine
func NewEngine(connection string) (*Engine, error) {
	return nil, nil
}

// SetLogger устанавливает логгер событий
func (self *Engine) SetLogger(logger log.ILogger) {
	self.logger = logger
}

// SetIncludedObjects устанавливает фильтр, позволяющий выбирать только те объекты БД, которые должны быть обработаны
func (self *Engine) SetIncludedObjects(filter filter.IFilter) {
	self.include = filter
}

// SetExcludedObjects устанавливает фильтр, позволяющий игнорировать объекты БД, которые должны быть заторонуты
// обработкой
func (self *Engine) SetExcludedObjects(filter filter.IFilter) {
	self.exclude = filter
}

// SetOutputDirectoryStructure устанавливает описание структуры каталога, где будут созданы скрипты
func (self *Engine) SetOutputDirectoryStructure(dirStruct dir.IStructure) {
	self.output = dirStruct
}
