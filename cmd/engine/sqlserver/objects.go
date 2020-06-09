package sqlserver

import (
	"database/sql"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

// IDatabaseObject базовый интерфейс объекта БД
type IDatabaseObject interface {
	// Catalog возвращает наименование базы данных
	Catalog() string
	// Schema возвращает схему базы данных
	Schema() string
	// Name возвращает наименование объекта БД
	Name() string
	// Definition возвращает определение объекта БД
	Definition() []byte
	// SetDefinition записывает новое определение объекта БД
	SetDefinition(data []byte)
	// DefinitionExists проверяет наличие определения у объекта БД
	DefinitionExists() bool
	// Type возвращает тип объекта БД
	Type() output.DatabaseObjectType
	// SchemaAndName возвращает наименование объекта в формате %schema%.%name%
	SchemaAndName() string
	// Owner возвращает владельца объекта БД
	Owner() string
}

// ISQLModule интерфейс SQL модуля (процедура, скалярная/табличная функция, представление, триггер...)
type ISQLModule interface {
	IDatabaseObject

	// ANSINulls объект создан с SET ANSI_NULLS ON
	ANSINulls() bool
	// ANSINullsValid сведения о ANSI_NULLS не равны NULL
	ANSINullsValid() bool
	// QuotedIdentifier объект создан с SET QUOTED_IDENTIFIER ON
	QuotedIdentifier() bool
	// QuotedIdentifierValid сведения о QUOTED_IDENTIFIER не равны NULL
	QuotedIdentifierValid() bool
}

// databaseObject базовая структура объекта БД
type databaseObject struct {
	// catalog наименование БД
	catalog sql.NullString
	// schema схема БД
	schema sql.NullString
	// name наименование объекта БД
	name sql.NullString
	// objectType тип объекта БД
	objectType sql.NullString
	// definition SQL код объекта БД
	definition sql.NullString
	// owner владелец объекта БД
	owner sql.NullString
}

// Catalog наименование базы данных
func (object databaseObject) Catalog() string {
	if object.catalog.Valid {
		return object.catalog.String
	}

	return ""
}

// Schema схема базы данных
func (object databaseObject) Schema() string {
	if object.schema.Valid {
		return object.schema.String
	}

	return ""
}

// Name наименование объекта БД
func (object databaseObject) Name() string {
	if object.name.Valid {
		return object.name.String
	}

	return ""
}

// Definition определение объекта БД
func (object databaseObject) Definition() []byte {
	if object.definition.Valid {
		return []byte(object.definition.String)
	}

	return nil
}

// SetDefinition записывает новое определение объекта БД
func (object *databaseObject) SetDefinition(data []byte) {
	object.definition = sql.NullString{
		String: string(data),
		Valid:  true,
	}
}

// DefinitionExists проверяет наличие определения у объекта БД
func (object databaseObject) DefinitionExists() bool {
	if !object.definition.Valid {
		return false
	}

	return len(object.definition.String) > 0
}

// Type тип объекта БД
func (object databaseObject) Type() output.DatabaseObjectType {
	if !object.objectType.Valid {
		return output.UnknownObject
	}

	objectType := object.objectType.String

	switch objectType {
	case "DATABASE":
		return output.Database
	case "SCHEMA":
		return output.Schema
	case "DOMAIN":
		return output.Domain
	case "BASE TABLE":
		return output.Table
	case "VIEW":
		return output.View
	case "TRIGGER":
		return output.Trigger
	case "FUNCTION":
		return output.Function
	case "PROCEDURE":
		return output.Procedure
	default:
		return output.UnknownObject
	}
}

// SchemaAndName наименование объекта в формате %schema%.%name%
func (object databaseObject) SchemaAndName() string {
	schema := object.Schema()
	name := object.Name()

	return SchemaAndObject(schema, name)
}

// Owner владелец объекта БД
func (object databaseObject) Owner() string {
	if object.owner.Valid {
		return object.owner.String
	}

	return ""
}

type module struct {
	databaseObject

	usesANSINulls        sql.NullBool
	usesQuotedIdentifier sql.NullBool
}

func (mod module) ANSINulls() bool {
	if mod.usesANSINulls.Valid {
		return mod.usesANSINulls.Bool
	}

	return false
}

func (mod module) ANSINullsValid() bool {
	return mod.usesANSINulls.Valid
}

func (mod module) QuotedIdentifier() bool {
	if mod.usesQuotedIdentifier.Valid {
		return mod.usesQuotedIdentifier.Bool
	}

	return false
}

func (mod module) QuotedIdentifierValid() bool {
	return mod.usesQuotedIdentifier.Valid
}
