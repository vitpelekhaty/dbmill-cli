package engine

import (
	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/commands"
	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/sqlserver"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

// IEngine интерфейс "движка" БД
type IEngine interface {
	// SetLogger устанавливает логгер событий
	SetLogger(logger log.ILogger)
	// ScriptsFolder создает скрипты объектов БД по указанному пути path
	ScriptsFolder(options ...commands.ScriptsFolderOption) commands.IScriptsFolderCommand
}

// Option опция "движка" базы данных
type Option func(engine IEngine)

// WithLogger указывает "движку" использовать логгер logger
func WithLogger(logger log.ILogger) Option {
	return func(engine IEngine) {
		engine.SetLogger(logger)
	}
}

// New возвращает экземпляр "движка" БД
func New(connection string, options ...Option) (IEngine, error) {
	engn, err := engine(connection)

	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(engn)
	}

	return engn, nil
}

func engine(connection string) (IEngine, error) {
	rdbms, err := RDBMS(connection)

	if err != nil {
		return nil, err
	}

	switch rdbms {
	case RDBMSSQLServer:
		return sqlserver.NewEngine(connection)
	default:
		return nil, ErrorUnsupportedDatabaseType
	}
}
