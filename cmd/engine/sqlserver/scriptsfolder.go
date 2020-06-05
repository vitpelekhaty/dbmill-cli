package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/reactivex/rxgo/v2"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/commands"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

// ScriptsFolderCommand реализация интерфейса IScriptsFolderCommand для SQL Server
type ScriptsFolderCommand struct {
	engine             *Engine
	include            filter.IFilter
	exclude            filter.IFilter
	decrypt            bool
	includeStaticData  bool
	types              map[output.DatabaseObjectType]bool
	definitionCallback commands.ObjectDefinitionCallback
}

// NewScriptsFolderCommand конструктор ScriptsFolderCommand
func NewScriptsFolderCommand(engine *Engine, options ...commands.ScriptsFolderOption) *ScriptsFolderCommand {
	command := &ScriptsFolderCommand{
		engine:             engine,
		include:            nil,
		exclude:            nil,
		decrypt:            false,
		includeStaticData:  false,
		types:              nil,
		definitionCallback: nil,
	}

	for _, option := range options {
		option(command)
	}

	return command
}

// Run запускает выполнение команды
func (command *ScriptsFolderCommand) Run() error {
	objects, err := command.databaseObjects()

	if err != nil {
		return err
	}

	in := command.enumObjects(objects)

	observable := rxgo.FromChannel(in)

	<-observable.
		Filter(func(item interface{}) bool {
			object := item.(databaseObject)
			return command.ObjectTypeIncluded(object.Type())
		}).
		Filter(func(item interface{}) bool {
			object := item.(databaseObject)
			return command.Included(object.SchemaAndName()) == nil
		}).
		Filter(func(item interface{}) bool {
			object := item.(databaseObject)
			return command.Excluded(object.SchemaAndName()) == filter.ErrorNotMatched
		}).
		Map(func(ctx context.Context, item interface{}) (interface{}, error) {
			obj := item.(databaseObject)
			return command.writeDefinition(ctx, obj)
		}).
		ForEach(func(item interface{}) {
			object := item.(databaseObject)
			command.engine.Log(log.DebugLevel, object.SchemaAndName())

			err := command.callObjectDefinitionCallback(object)

			if err != nil {
				command.engine.Log(log.ErrorLevel, err)
			}
		}, func(err error) {
			if err != nil {
				command.engine.Log(log.ErrorLevel, err)
			}
		}, func() {
			command.engine.Log(log.InfoLevel, "done")
		})

	return nil
}

func (command *ScriptsFolderCommand) enumObjects(databaseObjects []interface{}) chan rxgo.Item {
	out := make(chan rxgo.Item)

	go func() {
		defer close(out)

		for _, object := range databaseObjects {
			out <- rxgo.Of(object)
		}
	}()

	return out
}

func (command *ScriptsFolderCommand) callObjectDefinitionCallback(object databaseObject) error {
	if command.definitionCallback == nil {
		return nil
	}

	return command.definitionCallback(object.Catalog(), object.Schema(), object.Name(), object.Type(),
		object.Definition())
}

func (command *ScriptsFolderCommand) writeDefinition(ctx context.Context, object interface{}) (interface{}, error) {
	obj := object.(databaseObject)

	switch obj.Type() {
	case output.Schema:
		return command.writeSchemaDefinition(ctx, obj)
	}

	return object, nil
}

// SetIncludedObjects устанавливает фильтр, позволяющий выбирать только те объекты БД, которые должны быть
// обработаны
func (command *ScriptsFolderCommand) SetIncludedObjects(filter filter.IFilter) {
	command.include = filter
}

// SetExcludedObjects устанавливает фильтр, позволяющий игнорировать объекты БД, которые должны быть заторонуты
// обработкой
func (command *ScriptsFolderCommand) SetExcludedObjects(filter filter.IFilter) {
	command.exclude = filter
}

// SetObjectDefinitionCallback устанавливает callback для чтения определений объектов БД
func (command *ScriptsFolderCommand) SetObjectDefinitionCallback(callback commands.ObjectDefinitionCallback) {
	command.definitionCallback = callback
}

// StaticData опция выгрузки скриптов вставки данных
func (command *ScriptsFolderCommand) StaticData(on bool) {
	command.includeStaticData = on
}

// Decrypt по возможности расшифровывать определения объектов БД
func (command *ScriptsFolderCommand) Decrypt(on bool) {
	command.decrypt = on
}

// SetDatabaseObjectTypes устанавливает список типов объектов БД, которые необходимо выгрузить в скрипты
func (command *ScriptsFolderCommand) SetDatabaseObjectTypes(types []output.DatabaseObjectType) {
	t := make(map[output.DatabaseObjectType]bool)

	for _, objectType := range types {
		t[objectType] = true
	}

	command.types = t
}

// Included проверяет, должен ли объект object быть включен в обработку
func (command *ScriptsFolderCommand) Included(object string) error {
	if command.include == nil {
		return nil
	}

	return command.include.Match(object)
}

// Excluded проверяет, должен ли объект object быть исключен из обработки
func (command *ScriptsFolderCommand) Excluded(object string) error {
	if command.exclude == nil {
		return filter.ErrorNotMatched
	}

	return command.exclude.Match(object)
}

// ObjectTypeIncluded проверяет, должен ли тип объекта БД включен в обработку
func (command *ScriptsFolderCommand) ObjectTypeIncluded(object output.DatabaseObjectType) bool {
	if len(command.types) == 0 {
		return false
	}

	_, ok := command.types[object]

	return ok
}

func (command *ScriptsFolderCommand) databaseObjects() ([]interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	stmt, err := command.engine.db.PrepareContext(ctx, selectObjects)

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

		object := databaseObject{
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
func (object *databaseObject) DefinitionExists() bool {
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

	if strings.Trim(schema, " ") != "" {
		if strings.Trim(name, " ") != "" {
			return fmt.Sprintf("%s.%s", schema, name)
		}

		return schema
	} else {
		if strings.Trim(name, " ") != "" {
			return name
		}
	}

	return ""
}

// Owner владелец объекта БД
func (object databaseObject) Owner() string {
	if object.owner.Valid {
		return object.owner.String
	}

	return ""
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
