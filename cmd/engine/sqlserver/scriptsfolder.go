package sqlserver

import (
	"context"
	"database/sql"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/commands"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
)

// ScriptsFolderCommand реализация интерфейса IScriptsFolderCommand для SQL Server
type ScriptsFolderCommand struct {
	engine             *Engine
	include            filter.IFilter
	exclude            filter.IFilter
	decrypt            bool
	includeStaticData  bool
	definitionCallback commands.ObjectDefinitionCallback
}

// NewScriptsFolderCommand конструктор ScriptsFolderCommand
func NewScriptsFolderCommand(engine *Engine, options ...commands.ScriptFoldersOption) *ScriptsFolderCommand {
	command := &ScriptsFolderCommand{
		engine: engine,
	}

	for _, option := range options {
		option(command)
	}

	return command
}

func (self *ScriptsFolderCommand) Run() error {
	_, err := self.databaseObjects()

	if err != nil {
		return err
	}

	return nil
}

// SetIncludedObjects устанавливает фильтр, позволяющий выбирать только те объекты БД, которые должны быть
// обработаны
func (self *ScriptsFolderCommand) SetIncludedObjects(filter filter.IFilter) {
	self.include = filter
}

// SetExcludedObjects устанавливает фильтр, позволяющий игнорировать объекты БД, которые должны быть заторонуты
// обработкой
func (self *ScriptsFolderCommand) SetExcludedObjects(filter filter.IFilter) {
	self.exclude = filter
}

// SetObjectDefinitionCallback устанавливает callback для чтения определений объектов БД
func (self *ScriptsFolderCommand) SetObjectDefinitionCallback(callback commands.ObjectDefinitionCallback) {
	self.definitionCallback = callback
}

// StaticData опция выгрузки скриптов вставки данных
func (self *ScriptsFolderCommand) StaticData(on bool) {
	self.includeStaticData = on
}

// Decrypt по возможности расшифровывать определения объектов БД
func (self *ScriptsFolderCommand) Decrypt(on bool) {
	self.decrypt = on
}

// Included проверяет, должен ли объект object быть включен в обработку
func (self *ScriptsFolderCommand) Included(object string) error {
	if self.include == nil {
		return nil
	}

	return self.include.Match(object)
}

// Excluded проверяет, должен ли объект object быть исключен из обработки
func (self *ScriptsFolderCommand) Excluded(object string) error {
	if self.exclude == nil {
		return filter.ErrorNotMatched
	}

	return self.exclude.Match(object)
}

func (self *ScriptsFolderCommand) databaseObjects() ([]interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	stmt, err := self.engine.db.PrepareContext(ctx, selectObjects)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	objects := make([]interface{}, 0)

	var (
		catalog    sql.NullString
		schema     sql.NullString
		name       sql.NullString
		objectType sql.NullString
		definition sql.NullString
		owner      sql.NullString
	)

	for rows.Next() {
		err = rows.Scan(&catalog, &schema, &name, &objectType, &definition, &owner)

		if err != nil {
			return nil, err
		}

		object := &databaseObject{
			catalog:    catalog,
			schema:     schema,
			name:       name,
			objectType: objectType,
			definition: definition,
			owner:      owner,
		}

		objects = append(objects, object)
	}

	return objects, nil
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

const selectObjects = `
select info.catalog, info.[schema], info.name, info.type, info.definition,
       info.owner
from (
    select
           [order] = 1,
           [catalog] = schemas.CATALOG_NAME,
           [schema] = schemas.SCHEMA_NAME,
           [name] = schemas.SCHEMA_NAME,
           [type] = N'SCHEMA',
           [definition] = null,
           [owner] = isnull(users.name, schemas.SCHEMA_OWNER)
    from INFORMATION_SCHEMA.SCHEMATA as schemas
        inner join sys.schemas as ss on (schemas.SCHEMA_NAME = ss.name)
            inner join sys.sysusers as users on (ss.principal_id = users.uid) and (users.hasdbaccess != 0)
    union
    select
           [order] = 2,
           [catalog] = domains.DOMAIN_CATALOG,
           [schema] = isnull(schema_name(types.schema_id), domains.DOMAIN_SCHEMA),
           [name] = domains.DOMAIN_NAME,
           [type] = N'DOMAIN',
           [definition] = null,
           [owner] = null
    from INFORMATION_SCHEMA.DOMAINS as domains
        inner join sys.types as types on (domains.DOMAIN_NAME = types.name)
    union
    select
           [order] = 3,
           [catalog] = tables.TABLE_CATALOG,
           [schema] = isnull(schema_name(objects.schema_id), tables.TABLE_SCHEMA),
           [name] = tables.TABLE_NAME,
           [type] = tables.TABLE_TYPE,
           [definition] = null,
           [owner] = null
    from INFORMATION_SCHEMA.TABLES as tables
        inner join sys.objects as objects on (tables.TABLE_NAME = objects.name)
    where tables.TABLE_TYPE = N'BASE TABLE'
    union
    select
           [order] = 4,
           [catalog] = views.TABLE_CATALOG,
           [schema] = isnull(schema_name(objects.schema_id), views.TABLE_SCHEMA),
           [name] = views.TABLE_NAME,
           [type] = N'VIEW',
           [definition] = isnull(object_definition(objects.object_id), views.VIEW_DEFINITION),
           [owner] = null
    from INFORMATION_SCHEMA.VIEWS as views
        inner join sys.objects as objects on (views.TABLE_NAME = objects.name)
    union
    select
           [order] = 5,
           [catalog] = db_name(),
           [schema] = schema_name(objects.schema_id),
           [name] = objects.name,
           [type] = N'TRIGGER',
           [definition] = object_definition(objects.object_id),
           [owner] = null
    from sys.objects as objects
    where objects.type = 'TA'
    union
    select
           [order] = 5,
           [catalog] = db_name(),
           [schema] = schema_name(objects.schema_id),
           [name] = objects.name,
           [type] = N'TRIGGER',
           [definition] = object_definition(objects.object_id),
           [owner] = null
    from sys.objects as objects
    where objects.type = 'TR'
    union
    select
           [order] = 6,
           [catalog] = functions.ROUTINE_CATALOG,
           [schema] = isnull(schema_name(objects.schema_id), functions.ROUTINE_SCHEMA),
           [name] = functions.ROUTINE_NAME,
           [type] = functions.ROUTINE_TYPE,
           [definition] = isnull(object_definition(objects.object_id), functions.ROUTINE_DEFINITION),
           [owner] = null
    from INFORMATION_SCHEMA.ROUTINES as functions
        inner join sys.objects as objects on (functions.ROUTINE_NAME = objects.name)
    where functions.ROUTINE_TYPE = N'FUNCTION'
    union
    select
           [order] = 7,
           [catalog] = procedures.ROUTINE_CATALOG,
           [schema] = isnull(schema_name(objects.schema_id), procedures.ROUTINE_SCHEMA),
           [name] = procedures.ROUTINE_NAME,
           [type] = procedures.ROUTINE_TYPE,
           [definition] = isnull(object_definition(objects.object_id), procedures.ROUTINE_DEFINITION),
           [owner] = null
    from INFORMATION_SCHEMA.ROUTINES as procedures
        inner join sys.objects as objects on (procedures.ROUTINE_NAME = objects.name)
    where procedures.ROUTINE_TYPE = N'PROCEDURE'
) as info
order by info.catalog, info.[order], info.type, info.[schema], info.name
`
