package sqlserver

import (
	"context"
	"database/sql"
	"errors"
)

// ErrorInvalidConnectionObject ошибка "Объект соединения равен nil"
var ErrorInvalidConnectionObject = errors.New("connection object cannot be nil")

// serverVersion возвращает версию SQL Server
func serverVersion(db *sql.DB, ctx context.Context) (int, error) {
	if db == nil {
		return 0, ErrorInvalidConnectionObject
	}

	stmt, err := db.PrepareContext(ctx, selectServerVersion)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var ver int

	err = stmt.QueryRowContext(ctx).Scan(&ver)

	if err != nil {
		return 0, err
	}

	return ver, nil
}

const selectServerVersion = `
select isnull(
    try_cast(
        left(
            cast(serverproperty('productversion') as nvarchar),
            patindex('%.%', cast(serverproperty('productversion') as nvarchar)) - 1
        ) as int
    ),
0) as version
`
