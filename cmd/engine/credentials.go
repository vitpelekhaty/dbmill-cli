package engine

import (
	"net/url"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/sqlserver"
)

// Credentials возвращает имя пользователя и пароль, извлеченные из строки соединения с базой данных
func Credentials(connection string) (username, password string, err error) {
	u, err := url.Parse(connection)

	if err != nil {
		return "", "", err
	}

	switch u.Scheme {
	case "sqlserver":
		return sqlserver.Credentials(connection)
	default:
		return "", "", ErrorUnsupportedDatabaseType
	}
}
