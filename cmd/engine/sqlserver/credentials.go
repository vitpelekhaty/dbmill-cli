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

// SetCredentials заменяет имя пользователя и пароль в строке соединения с БД
func SetCredentials(connection, username, password string) (string, error) {
	u, err := url.Parse(connection)

	if err != nil {
		return connection, err
	}

	emptyUsername := strings.Trim(username, " ") == ""
	emptyPassword := strings.Trim(password, " ") == ""

	user := u.User

	if user.Username() != "" {
		if !emptyPassword {
			u.User = url.UserPassword(username, password)
		} else {
			if !emptyUsername {
				u.User = url.User(username)
			} else {
				u.User = nil
			}
		}

		return u.String(), nil
	}

	q, err := url.ParseQuery(u.RawQuery)

	if err != nil {
		return connection, err
	}

	un := q.Get("user id")

	if user.Username() == "" && un == "" {
		if !emptyPassword {
			u.User = url.UserPassword(username, password)
		} else {
			if !emptyUsername {
				u.User = url.User(username)
			} else {
				u.User = nil
			}
		}

		return u.String(), nil
	}

	if !emptyUsername {
		q.Set("user id", username)
	} else {
		q.Del("user id")
	}

	if !emptyPassword {
		q.Set("password", password)
	} else {
		q.Del("password")
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}
