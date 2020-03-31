package engine

import (
	"errors"
	"net/url"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/sqlserver"
)

var ErrorUnsupportedDatabaseType = errors.New("unsupported database type")

// NewDatabaseConnection новое соединение с базой данных
func NewDatabaseConnection(connection string, options ...DatabaseOption) (IDatabase, error) {
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
