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
	indexes          ObjectsIndexes
	foreignKeys      ObjectsForeignKeys
	tables           Tables

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
		foreignKeys:      nil,
		tables:           nil,

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
	case output.Table:
		return command.writeTableDefinition(ctx, obj)
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

	foreignKeys, err := command.metaReader.ForeignKeys(ctx)

	if err != nil {
		return err
	}

	command.foreignKeys = foreignKeys

	tables, err := command.metaReader.Tables(ctx)

	if err != nil {
		return err
	}

	command.tables = tables

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
			description          sql.NullString
		)

		var objType string

		for rows.Next() {
			err = rows.Scan(&catalog, &schema, &name, &objectType, &definition, &owner, &usesANSINulls,
				&usesQuotedIdentifier, &description)

			if err == nil {
				if objectType.Valid {
					objType = objectType.String
				}

				var object interface{}

				switch objType {
				case "FUNCTION", "PROCEDURE", "TRIGGER", "VIEW":
					object = &module{
						databaseObject: databaseObject{
							catalog:     catalog,
							schema:      schema,
							name:        name,
							objectType:  objectType,
							definition:  definition,
							owner:       owner,
							description: description,
						},
						usesANSINulls:        usesANSINulls,
						usesQuotedIdentifier: usesQuotedIdentifier,
					}
				default:
					object = &databaseObject{
						catalog:     catalog,
						schema:      schema,
						name:        name,
						objectType:  objectType,
						definition:  definition,
						owner:       owner,
						description: description,
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
with tableTypes (object_id, schema_id, name) as (
    select table_types.type_table_object_id as object_id, types.schema_id, types.name
    from sys.types as types
        inner join sys.table_types as table_types on (types.user_type_id = table_types.user_type_id)
    where (types.is_table_type != cast(0 as bit)) and (types.is_user_defined != cast(0 as bit))
),
extendedProperties (object_id, description, class) as (
    select props.major_id as object_id, cast(props.value as nvarchar(2048)) as description, props.class
    from sys.extended_properties as props
    where (name = N'MS_Description') and (props.minor_id = 0)
),
objectDescriptions (object_id, description, class) as (
    select props.object_id,  props.description, props.class
    from extendedProperties as props
    where (props.class = 1)
    union
    select props.object_id,  props.description, props.class
    from extendedProperties as props
    where (props.class = 3)
    union
    select props.object_id,  props.description, props.class
    from extendedProperties as props
    where (props.class = 6)
)
select info.catalog, info.[schema], info.name, info.type, info.definition,
       info.owner, info.uses_quoted_identifier, info.uses_ansi_nulls, info.description
from (
    select
        [order] = 1,
        [catalog] = db_name(),
        [schema] = schemas.name,
        [name] = null,
        [type] = N'SCHEMA',
        [definition] = null,
        [owner] = users.name,
        [uses_ansi_nulls] = null,
        [uses_quoted_identifier] = null,
        [description] = prop.description
    from sys.schemas as schemas
        inner join sys.sysusers as users on (schemas.principal_id = users.uid) and (users.hasdbaccess != 0)
        left join objectDescriptions as prop on (schemas.schema_id = prop.object_id) and (prop.class = 3)
    union
    select
        [order] = 2,
        [catalog] = db_name(),
        [schema] = schema_name(types.schema_id),
        [name] = types.name,
        [type] = N'DATA TYPE',
        [definition] = null,
        [owner] = null,
        [uses_ansi_nulls] = null,
        [uses_quoted_identifier] = null,
        [description] = prop.description
    from sys.types as types
        left join objectDescriptions as prop on (types.user_type_id = prop.object_id) and (prop.class = 6)
    where (types.is_user_defined != cast(0 as bit)) and (types.is_table_type = cast(0 as bit))
    union
    select
        [order] = case objects.type
            when 'TT' then 2
            when 'U' then 3
            when 'V' then 4
            when 'TR' then 5
            when 'FN' then 6
            when 'IF' then 6
            when 'TF' then 6
            when 'P' then 7
            else null
        end,

        [catalog] = db_name(),
        [schema] = schema_name(iif(objects.type = 'TT', tableTypes.schema_id, objects.schema_id)),
        [name] = iif(objects.type = 'TT', tableTypes.name, objects.name),

        [type] = case objects.type
            when 'TT' then N'TABLE TYPE'
            when 'U' then N'BASE TABLE'
            when 'V' then N'VIEW'
            when 'TR' then N'TRIGGER'
            when 'FN' then N'FUNCTION'
            when 'IF' then N'FUNCTION'
            when 'TF' then N'FUNCTION'
            when 'P' then N'PROCEDURE'
            else null
        end,

        [definition] = object_definition(objects.object_id),
        [owner] = null,
        [uses_ansi_nulls] = modules.uses_ansi_nulls,
        [uses_quoted_identifier] = modules.uses_quoted_identifier,
        [description] = iif(objects.type = 'TT', prop_types.description, prop_objects.description)

    from sys.objects as objects
        left join sys.sql_modules as modules on (objects.object_id = modules.object_id)
        left join tableTypes on (objects.object_id = tableTypes.object_id)
        left join objectDescriptions as prop_objects on (objects.object_id = prop_objects.object_id)
            and (prop_objects.class = 1)
        left join objectDescriptions as prop_types on (objects.object_id = prop_types.object_id)
            and (prop_types.class = 6)
    where objects.type in ('TT', 'U', 'V', 'TR', 'FN', 'IF', 'TF', 'P')
) as info
order by info.catalog, info.[order], info.type, info.[schema], info.name
`
