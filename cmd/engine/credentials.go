package engine

import (
	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/sqlserver"
)

// Credentials возвращает имя пользователя и пароль, извлеченные из строки соединения с базой данных
func Credentials(connection string) (username, password string, err error) {
	rdbms, err := RDBMS(connection)

	if err != nil {
		return "", "", err
	}

	switch rdbms {
	case RDBMSSQLServer:
		return sqlserver.Credentials(connection)
	default:
		return "", "", ErrorUnsupportedDatabaseType
	}
}

// SetCredentials заменяет имя пользователя и пароль в строке соединения с БД
func SetCredentials(connection, username, password string) (string, error) {
	rdbms, err := RDBMS(connection)

	if err != nil {
		return "", err
	}

	switch rdbms {
	case RDBMSSQLServer:
		return sqlserver.SetCredentials(connection, username, password)
	default:
		return "", ErrorUnsupportedDatabaseType
	}
}
