package engine

import (
	"net/url"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/sqlserver"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/dir"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

// DatabaseOption опция "движка" базы данных
type DatabaseOption func(engine IDatabase)

// WithLogger указывает "движку" использовать логгер logger
func WithLogger(logger log.ILogger) DatabaseOption {
	return func(engine IDatabase) {
		engine.SetLogger(logger)
	}
}

// WithIncludedObjects указывает "движку" использовать фильтр filter, чтобы ограничить список объектов БД указанным
// пользователем
func WithIncludedObjects(filter filter.IFilter) DatabaseOption {
	return func(engine IDatabase) {
		engine.SetIncludedObjects(filter)
	}
}

// WithExcludedObjects указывает "движку" использовать фильтр filter, чтобы исключить из обработки указанные
// пользователем объекты БД
func WithExcludedObjects(filter filter.IFilter) DatabaseOption {
	return func(engine IDatabase) {
		engine.SetExcludedObjects(filter)
	}
}

// WithOutputDirStructure указывает "движку" при создании скриптов руководствоваться указанной структурой целевого
// каталога
func WithOutputDirStructure(dirStruct dir.IStructure) DatabaseOption {
	return func(engine IDatabase) {
		engine.SetOutputDirectoryStructure(dirStruct)
	}
}

// NewDatabaseConnection новое соединение с базой данных
func NewDatabaseConnection(connection string, options ...DatabaseOption) (IDatabase, error) {
	db, err := engine(connection)

	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(db)
	}

	return db, nil
}

func engine(connection string) (IDatabase, error) {
	u, err := url.Parse(connection)

	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "sqlserver":
		return sqlserver.NewEngine(connection)
	default:
		return nil, ErrorUnsupportedDatabaseType
	}
}
