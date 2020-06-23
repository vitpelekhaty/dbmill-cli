package sqlserver

// Database параметры базы данных
type Database struct {
}

// String возвращает скрипт создания базы данных
func (db *Database) String() string {
	return ""
}

const selectDatabaseCollation = `select isnull(DATABASEPROPERTYEX(db_name(), 'Collation'), N'') AS collation`
