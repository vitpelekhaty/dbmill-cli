package engine

import (
	"errors"
	"net/url"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/sqlserver"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

type DatabaseOption func(engine IDatabase)

// WithLogger устанавливает
func WithLogger(logger *log.Logger) DatabaseOption {
	return func(engine IDatabase) {
		engine.SetLogger(logger)
	}
}

// ErrorUnsupportedDatabaseType ошибка "Неподдерживаемая СУБД"
var ErrorUnsupportedDatabaseType = errors.New("unsupported database type")

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
