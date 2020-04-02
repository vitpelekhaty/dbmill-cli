package sqlserver

import (
	"context"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"

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
	db, err := sql.Open("sqlserver", connection)

	if err != nil {
		return nil, err
	}

	err = db.PingContext(context.Background())

	if err != nil {
		return nil, err
	}

	return &Engine{
		db:      db,
		logger:  nil,
		include: nil,
		exclude: nil,
		output:  dir.Default,
	}, nil
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
	if dirStruct == nil {
		self.output = dir.Default
	} else {
		self.output = dirStruct
	}
}

// Log создает запись в логе, если указан логгер
func (self *Engine) Log(level log.Level, args ...interface{}) {
	if self.logger == nil {
		return
	}

	self.logger.Print(level, args...)
}

// Logf создает фомрматированную запись в логе, если указан логгер
func (self *Engine) Logf(level log.Level, format string, args ...interface{}) {
	if self.logger == nil {
		return
	}

	self.logger.Printf(level, format, args...)
}

// Included проверяет, должен ли объект object быть включен в обработку
func (self *Engine) Included(object string) error {
	if self.include == nil {
		return nil
	}

	return self.include.Match(object)
}

// Excluded проверяет, должен ли объект object быть исключен из обработки
func (self *Engine) Excluded(object string) error {
	if self.exclude == nil {
		return filter.ErrorNotMatched
	}

	return self.exclude.Match(object)
}

// OutputDirectoryItemInfo возвращает целевой каталог и маску имени файла для указанного типа объекта itemType.
// Если информация не найдена, то в параметре ok возвращается false, в противном случае - true
func (self *Engine) OutputDirectoryItemInfo(item dir.StructItemType) (subdirectory, mask string, ok bool) {
	return self.output.ItemInfo(item)
}
