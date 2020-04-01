package dir

// Default структура каталога скриптов по умолчанию
var Default *Structure

func init() {
	Default, _ = NewStructure([]byte(defaultData))
}

const defaultData = `
Database:
  subdirectory: Database/
  mask: "$object$.sql"

Table:
  subdirectory: Tables/
  mask: $schema$.$object$.sql

TableData:
  subdirectory: Tables/StaticData
  mask: $schema$.$object$.Data.sql
`
