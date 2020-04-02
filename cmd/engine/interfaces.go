package engine

import (
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/dir"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

// IDatabase интерфейс базы данных
type IDatabase interface {
	// SetLogger устанавливает логгер событий
	SetLogger(logger log.ILogger)
	// SetIncludedObjects устанавливает фильтр, позволяющий выбирать только те объекты БД, которые должны быть
	// обработаны
	SetIncludedObjects(filter filter.IFilter)
	// SetExcludedObjects устанавливает фильтр, позволяющий игнорировать объекты БД, которые должны быть заторонуты
	// обработкой
	SetExcludedObjects(filter filter.IFilter)
	// SetOutputDirectoryStructure устанавливает описание структуры каталога, где будут созданы скрипты
	SetOutputDirectoryStructure(dirStruct dir.IStructure)
	// ScriptsFolder создает скрипты объектов БД по указанному пути path
	ScriptsFolder(path string, includeData, decrypt bool) error
}
