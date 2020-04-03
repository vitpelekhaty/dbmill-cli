package sqlserver

import (
	"net/url"
	"strings"
)

// Credentials возвращает имя пользователя и пароль, извлеченные из строки соединения с базой данных
func Credentials(connection string) (username, password string, err error) {
	u, err := url.Parse(connection)

	if err != nil {
		return username, password, err
	}

	user := u.User

	username = user.Username()
	password, _ = user.Password()

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		q, err := url.ParseQuery(u.RawQuery)

		if err != nil {
			return username, password, err
		}

		if strings.Trim(username, " ") == "" {
			username = q.Get("user id")
		}

		if strings.Trim(password, " ") == "" {
			password = q.Get("password")
		}
	}

	return username, password, nil
}
