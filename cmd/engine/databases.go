package engine

import (
	"net/url"
)

// RDBMSType известные СУБД
type RDBMSType byte

const (
	// RDBMSUnknown неизвестная СУБД (ошибка)
	RDBMSUnknown RDBMSType = iota
	// RDBMSSQLServer SQL Server
	RDBMSSQLServer
)

// RDBMS возвращает тип СУБД, с которым предстоит работать
func RDBMS(connection string) (RDBMSType, error) {
	u, err := url.Parse(connection)

	if err != nil {
		return RDBMSUnknown, err
	}

	switch u.Scheme {
	case "sqlserver":
		return RDBMSSQLServer, nil
	default:
		return RDBMSUnknown, ErrorUnsupportedDatabaseType
	}
}
