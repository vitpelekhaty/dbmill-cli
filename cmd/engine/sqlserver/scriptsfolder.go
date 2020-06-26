package sqlserver

import (
	"context"
	"database/sql"

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
	skipPermissions    bool
	types              map[output.DatabaseObjectType]bool
	definitionCallback commands.ObjectDefinitionCallback
	metaReader         *MetadataReader

	userDefinedTypes UserDefinedTypes
	permissions      ObjectPermissions
	columns          ObjectColumns
	indexes          Indexes

	databaseCollation string
}

// NewScriptsFolderCommand конструктор ScriptsFolderCommand
func NewScriptsFolderCommand(engine *Engine, options ...commands.ScriptsFolderOption) *ScriptsFolderCommand {
	metaReader, _ := engine.MetadataReader()

	command := &ScriptsFolderCommand{
		engine:             engine,
		include:            nil,
		exclude:            nil,
		decrypt:            false,
		includeStaticData:  false,
		skipPermissions:    false,
		types:              nil,
		definitionCallback: nil,
		metaReader:         metaReader,

		permissions:      nil,
		userDefinedTypes: nil,
		columns:          nil,
		indexes:          nil,

		databaseCollation: "",
	}

	for _, option := range options {
		option(command)
	}

	return command
}

// Run запускает выполнение команды
func (command *ScriptsFolderCommand) Run() error {
	ctx := context.Background()

	command.engine.Log(log.DebugLevel, "metadata reading...")

	err := command.ReadMetadata(ctx)

	if err != nil {
		return err
	}

	objects, err := command.databaseObjects(ctx)

	if err != nil {
		return err
	}

	observable := rxgo.FromChannel(objects)

	<-observable.
		Filter(func(item interface{}) bool {
			object := item.(IDatabaseObject)
			return command.ObjectTypeIncluded(object.Type())
		}).
		Filter(func(item interface{}) bool {
			object := item.(IDatabaseObject)
			return command.Included(object.SchemaAndName(true)) == nil
		}).
		Filter(func(item interface{}) bool {
			object := item.(IDatabaseObject)
			return command.Excluded(object.SchemaAndName(true)) == filter.ErrorNotMatched
		}).
		Map(func(ctx context.Context, item interface{}) (interface{}, error) {
			obj := item.(IDatabaseObject)
			return command.writeDefinition(ctx, obj)
		}).
		ForEach(func(item interface{}) {
			object := item.(IDatabaseObject)
			command.engine.Log(log.DebugLevel, object.SchemaAndName(true))

			err := command.callObjectDefinitionCallback(object)

			if err != nil {
				command.engine.Log(log.ErrorLevel, err)
			}
		}, func(err error) {
			if err != nil {
				command.engine.Log(log.ErrorLevel, err)
			}
		}, func() {
			command.engine.Log(log.DebugLevel, "done")
		})

	return nil
}

func (command *ScriptsFolderCommand) callObjectDefinitionCallback(object IDatabaseObject) error {
	if command.definitionCallback == nil {
		return nil
	}

	return command.definitionCallback(object.Catalog(), object.Schema(), object.Name(), object.Type(),
		object.Definition())
}

func (command *ScriptsFolderCommand) writeDefinition(ctx context.Context, object interface{}) (interface{}, error) {
	obj := object.(IDatabaseObject)

	switch obj.Type() {
	case output.Schema:
		return command.writeSchemaDefinition(ctx, obj)
	case output.Procedure:
		return command.writeProcedureDefinition(ctx, obj)
	case output.Function:
		return command.writeFunctionDefinition(ctx, obj)
	case output.View:
		return command.writeViewDefinition(ctx, obj)
	case output.Trigger:
		return command.writeTriggerDefinition(ctx, obj)
	case output.UserDefinedTableType, output.UserDefinedDataType:
		return command.writeDomainDefinition(ctx, obj)
	}

	return object, nil
}

func (command *ScriptsFolderCommand) ReadMetadata(ctx context.Context) error {
	collation, err := command.metaReader.DatabaseCollation(ctx)

	if err != nil {
		return err
	}

	command.databaseCollation = collation

	permissions, err := command.metaReader.Permissions(ctx)

	if err != nil {
		return err
	}

	command.permissions = permissions

	userTypes, err := command.metaReader.UserDefinedTypes(ctx)

	if err != nil {
		return err
	}

	command.userDefinedTypes = userTypes

	columns, err := command.metaReader.ObjectColumns(ctx)

	if err != nil {
		return err
	}

	command.columns = columns

	indexes, err := command.metaReader.Indexes(ctx)

	if err != nil {
		return err
	}

	command.indexes = indexes

	return nil
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

// SkipPermissions не добавлять в скрипты разрешения на объект
func (command *ScriptsFolderCommand) SkipPermissions(on bool) {
	command.skipPermissions = on
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

// DatabaseCollation возвращает collation базы
func (command *ScriptsFolderCommand) DatabaseCollation() string {
	return command.databaseCollation
}

func (command *ScriptsFolderCommand) databaseObjects(ctx context.Context) (chan rxgo.Item, error) {
	stmt, err := command.engine.db.PrepareContext(ctx, selectObjects)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	out := make(chan rxgo.Item)

	go func() {
		defer close(out)
		defer rows.Close()
		defer stmt.Close()

		var (
			catalog              sql.NullString
			schema               sql.NullString
			name                 sql.NullString
			objectType           sql.NullString
			definition           sql.NullString
			owner                sql.NullString
			usesANSINulls        sql.NullBool
			usesQuotedIdentifier sql.NullBool
		)

		var objType string

		for rows.Next() {
			err = rows.Scan(&catalog, &schema, &name, &objectType, &definition, &owner, &usesANSINulls,
				&usesQuotedIdentifier)

			if err == nil {
				if objectType.Valid {
					objType = objectType.String
				}

				var object interface{}

				switch objType {
				case "FUNCTION", "PROCEDURE", "TRIGGER", "VIEW":
					object = &module{
						databaseObject: databaseObject{
							catalog:    catalog,
							schema:     schema,
							name:       name,
							objectType: objectType,
							definition: definition,
							owner:      owner,
						},
						usesANSINulls:        usesANSINulls,
						usesQuotedIdentifier: usesQuotedIdentifier,
					}
				default:
					object = &databaseObject{
						catalog:    catalog,
						schema:     schema,
						name:       name,
						objectType: objectType,
						definition: definition,
						owner:      owner,
					}
				}

				out <- rxgo.Of(object)
			} else {
				out <- rxgo.Error(err)
			}
		}
	}()

	return out, nil
}

const selectObjects = `
select info.Catalog, info.[Schema], info.name, info.type, info.definition,
       info.owner, info.uses_quoted_identifier, info.uses_ansi_nulls
from (
    select
        [order] = 1,
        [Catalog] = schemas.CATALOG_NAME,
        [Schema] = schemas.SCHEMA_NAME,
        [name] = null,
        [type] = N'SCHEMA',
        [definition] = null,
        [owner] = isnull(users.name, schemas.SCHEMA_OWNER),
        [uses_ansi_nulls] = null,
        [uses_quoted_identifier] = null
    from INFORMATION_SCHEMA.SCHEMATA as schemas
        inner join sys.schemas as ss on (schemas.SCHEMA_NAME = ss.name)
            inner join sys.sysusers as users on (ss.principal_id = users.uid) and (users.hasdbaccess != 0)
    union
    select
        [order] = 2,
        [Catalog] = db_name(),
        [Schema] = schema_name(types.schema_id),
        [name] = types.name,
        [type] = N'DATA TYPE',
        [definition] = null,
        [owner] = null,
        [uses_ansi_nulls] = null,
        [uses_quoted_identifier] = null
    from sys.types as types
    where (types.is_user_defined != cast(0 as bit)) and (types.is_table_type = cast(0 as bit))
    union
    select
        [order] = 2,
        [Catalog] = db_name(),
        [Schema] = schema_name(types.schema_id),
        [name] = types.name,
        [type] = N'TABLE TYPE',
        [definition] = null,
        [owner] = null,
        [uses_ansi_nulls] = null,
        [uses_quoted_identifier] = null
    from sys.types as types
    where (types.is_user_defined != cast(0 as bit)) and (types.is_table_type != cast(0 as bit))
    union
    select
        [order] = 3,
        [Catalog] = tables.TABLE_CATALOG,
        [Schema] = isnull(schema_name(objects.schema_id), tables.TABLE_SCHEMA),
        [name] = tables.TABLE_NAME,
        [type] = tables.TABLE_TYPE,
        [definition] = null,
        [owner] = null,
        [uses_ansi_nulls] = null,
        [uses_quoted_identifier] = null
    from INFORMATION_SCHEMA.TABLES as tables
        inner join sys.objects as objects on (tables.TABLE_NAME = objects.name)
    where tables.TABLE_TYPE = N'BASE TABLE'
    union
    select
        [order] = 4,
        [Catalog] = views.TABLE_CATALOG,
        [Schema] = isnull(schema_name(objects.schema_id), views.TABLE_SCHEMA),
        [name] = views.TABLE_NAME,
        [type] = N'VIEW',
        [definition] = isnull(object_definition(objects.object_id), views.VIEW_DEFINITION),
        [owner] = null,
        [uses_ansi_nulls] = modules.uses_ansi_nulls,
        [uses_quoted_identifier] = modules.uses_quoted_identifier
    from INFORMATION_SCHEMA.VIEWS as views
        inner join sys.objects as objects on (views.TABLE_NAME = objects.name)
            inner join sys.sql_modules as modules on (objects.object_id = modules.object_id)
    union
    select
        [order] = 5,
        [Catalog] = db_name(),
        [Schema] = schema_name(objects.schema_id),
        [name] = objects.name,
        [type] = N'TRIGGER',
        [definition] = object_definition(objects.object_id),
        [owner] = null,
        [uses_ansi_nulls] = modules.uses_ansi_nulls,
        [uses_quoted_identifier] = modules.uses_quoted_identifier
    from sys.objects as objects
        inner join sys.sql_modules as modules on (objects.object_id = modules.object_id)
    where objects.type = 'TR'
    union
    select
        [order] = 6,
        [Catalog] = functions.ROUTINE_CATALOG,
        [Schema] = isnull(schema_name(objects.schema_id), functions.ROUTINE_SCHEMA),
        [name] = functions.ROUTINE_NAME,
        [type] = functions.ROUTINE_TYPE,
        [definition] = isnull(object_definition(objects.object_id), functions.ROUTINE_DEFINITION),
        [owner] = null,
        [uses_ansi_nulls] = modules.uses_ansi_nulls,
        [uses_quoted_identifier] = modules.uses_quoted_identifier
    from INFORMATION_SCHEMA.ROUTINES as functions
        inner join sys.objects as objects on (functions.ROUTINE_NAME = objects.name)
            inner join sys.sql_modules as modules on (objects.object_id = modules.object_id)
    where functions.ROUTINE_TYPE = N'FUNCTION'
    union
    select
        [order] = 7,
        [Catalog] = procedures.ROUTINE_CATALOG,
        [Schema] = isnull(schema_name(objects.schema_id), procedures.ROUTINE_SCHEMA),
        [name] = procedures.ROUTINE_NAME,
        [type] = procedures.ROUTINE_TYPE,
        [definition] = isnull(object_definition(objects.object_id), procedures.ROUTINE_DEFINITION),
        [owner] = null,
        [uses_ansi_nulls] = modules.uses_ansi_nulls,
        [uses_quoted_identifier] = modules.uses_quoted_identifier
    from INFORMATION_SCHEMA.ROUTINES as procedures
        inner join sys.objects as objects on (procedures.ROUTINE_NAME = objects.name)
            inner join sys.sql_modules as modules on (objects.object_id = modules.object_id)
    where procedures.ROUTINE_TYPE = N'PROCEDURE'
) as info
order by info.Catalog, info.[order], info.type, info.[Schema], info.name
`
