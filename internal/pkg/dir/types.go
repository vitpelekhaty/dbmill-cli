package dir

// StructItemType тип элемента структуры каталога
type StructItemType uint16

const (
	// Database база данных
	Database StructItemType = iota
	// Table таблица
	Table
	// TableData данные таблицы
	TableData
	// View представление
	View
	// Procedure процедура
	Procedure
	// ScalarFunc скалярная функция
	ScalarFunc
	// TableFunc табличная функция
	TableFunc
)

var itemTypeMapping = map[string]StructItemType{
	"Database":   Database,
	"Table":      Table,
	"TableData":  TableData,
	"View":       View,
	"Procedure":  Procedure,
	"ScalarFunc": ScalarFunc,
	"TableFunc":  TableFunc,
}
