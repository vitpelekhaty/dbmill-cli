package output

// DatabaseObjectType тип элемента структуры каталога
type DatabaseObjectType uint16

const (
	// UnknownObject неизвестный тип объекта БД
	UnknownObject DatabaseObjectType = iota
	// Database база данных
	Database
	// Table таблица
	Table
	// StaticData данные таблицы
	StaticData
	// View представление
	View
	// Procedure процедура
	Procedure
	// Function скалярная функция
	Function
	// Trigger табличная функция
	Trigger
	// Domain пользовательский тип данных
	Domain
	// Schema
	Schema
)

// String возвращает строковое представление значения типа DatabaseObjectType
func (objType DatabaseObjectType) String() string {
	return databaseObjectTypeMapping[objType]
}

var databaseObjectTypeMapping = map[DatabaseObjectType]string{
	UnknownObject: "unknown",
	Database:      "database",
	Table:         "table",
	StaticData:    "staticData",
	View:          "view",
	Procedure:     "procedure",
	Function:      "function",
	Trigger:       "trigger",
	Domain:        "domain",
	Schema:        "schema",
}

var databaseObjectTypeMappingReverse = map[string]DatabaseObjectType{
	"unknown":    UnknownObject,
	"database":   Database,
	"table":      Table,
	"staticData": StaticData,
	"view":       View,
	"procedure":  Procedure,
	"function":   Function,
	"trigger":    Trigger,
	"domain":     Domain,
	"schema":     Schema,
}
