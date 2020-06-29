package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

type UserDefinedType struct {
	Catalog           string
	TypeName          string
	Schema            string
	parentTypeName    sql.NullString
	maxLength         sql.NullString
	precision         sql.NullInt32
	scale             sql.NullInt32
	collation         sql.NullString
	IsNullable        bool
	IsTableType       bool
	IsMemoryOptimized bool
}

// SchemaAndObject возвращает наименование объекта в формате %Schema%.%name%
func (udt UserDefinedType) SchemaAndName(useBrackets bool) string {
	return SchemaAndObject(udt.Schema, udt.TypeName, useBrackets)
}

// HasMaxLength проверяет, указана ли максимальная длина (в байтах) для типа
func (udt UserDefinedType) HasMaxLength() bool {
	return udt.maxLength.Valid
}

// MaxLength максимальная длина (в байтах) типа
func (udt UserDefinedType) MaxLength() string {
	if udt.maxLength.Valid {
		return udt.maxLength.String
	}

	return ""
}

// HasPrecision проверяет, указана ли точность
func (udt UserDefinedType) HasPrecision() bool {
	return udt.precision.Valid
}

// Precision точность
func (udt UserDefinedType) Precision() int {
	if udt.precision.Valid {
		return int(udt.precision.Int32)
	}

	return 0
}

// HasScale проверяет, указан ли масштаб
func (udt UserDefinedType) HasScale() bool {
	return udt.scale.Valid
}

// Scale масштаб
func (udt UserDefinedType) Scale() int {
	if udt.scale.Valid {
		return int(udt.scale.Int32)
	}

	return 0
}

// HasCollation провеяет, указан ли collation
func (udt UserDefinedType) HasCollation() bool {
	return udt.collation.Valid
}

// Collation
func (udt UserDefinedType) Collation() string {
	if udt.collation.Valid {
		return udt.collation.String
	}

	return ""
}

func (udt UserDefinedType) ParentTypeName() string {
	if udt.parentTypeName.Valid && !udt.IsTableType {
		return udt.parentTypeName.String
	}

	return ""
}

// ParentType возвращает полное описание родительского типа
func (udt UserDefinedType) ParentType() string {
	var builder strings.Builder

	builder.WriteString("[" + udt.ParentTypeName() + "]")

	openedParentheses := udt.HasMaxLength() || udt.HasPrecision()

	if openedParentheses {
		builder.WriteRune('(')
	}

	if udt.HasMaxLength() {
		builder.WriteString(udt.MaxLength())
	} else {
		if udt.HasPrecision() {
			precision := udt.Precision()
			scale := udt.Scale()

			builder.WriteString(strconv.Itoa(precision))

			if scale > 0 {
				builder.WriteString(", ")
				builder.WriteString(strconv.Itoa(scale))
			}
		}
	}

	if openedParentheses {
		builder.WriteRune(')')
	}

	return builder.String()
}

type UserDefinedTypes map[string]*UserDefinedType

func (types UserDefinedTypes) append(udt *UserDefinedType) {
	if udt == nil {
		return
	}

	types[udt.SchemaAndName(true)] = udt
}

func (command *ScriptsFolderCommand) writeDomainDefinition(ctx context.Context, object interface{}) (interface{},
	error) {
	obj, ok := object.(IDatabaseObject)

	if !ok {
		return object, errors.New("object is not a database object")
	}

	name := obj.SchemaAndName(true)

	if obj.Type() != output.UserDefinedDataType && obj.Type() != output.UserDefinedTableType {
		return object, fmt.Errorf("object %s is not a domain", name)
	}

	domain, ok := command.userDefinedTypes[name]

	if !ok {
		return object, fmt.Errorf("no info about domain %s", name)
	}

	if obj.Type() == output.UserDefinedDataType {
		if !domain.IsTableType {
			return command.writeDataTypeDefinition(ctx, obj, domain)
		} else {
			return obj, fmt.Errorf("%s is not data type", name)
		}
	}

	if obj.Type() == output.UserDefinedTableType {
		if domain.IsTableType {
			return command.writeTableTypeDefinition(ctx, obj, domain)
		} else {
			return obj, fmt.Errorf("%s is not table type", name)
		}
	}

	return obj, nil
}

func (command *ScriptsFolderCommand) writeDataTypeDefinition(ctx context.Context, object IDatabaseObject,
	domain *UserDefinedType) (IDatabaseObject, error) {
	const dataTypeDefinition = "CREATE TYPE %s FROM %s\nGO"

	var builder strings.Builder

	userTypeName := domain.SchemaAndName(true)

	builder.WriteString(domain.ParentType())

	if !domain.IsNullable {
		builder.WriteString(" NOT NULL")
	}

	definition := fmt.Sprintf(dataTypeDefinition, userTypeName, builder.String())
	object.SetDefinition([]byte(definition))

	return object, nil
}

func (command *ScriptsFolderCommand) writeTableTypeDefinition(ctx context.Context, object IDatabaseObject,
	domain *UserDefinedType) (IDatabaseObject, error) {
	var builder strings.Builder
	userTypeName := domain.SchemaAndName(true)

	builder.WriteString(fmt.Sprintf("CREATE TYPE %s AS TABLE", userTypeName))

	columns := command.columns[userTypeName]

	if len(columns) > 0 {
		builder.WriteString(" (")

		cols := columns.List()
		sort.Slice(cols, func(i, j int) bool {
			return cols[i].ID < cols[j].ID
		})

		for index, col := range cols {
			col.SetOptions(WithColumnOwner(OwnerUserDefinedTableDataType),
				WithDefaultCollation(command.DatabaseCollation()))

			if index > 0 {
				builder.WriteRune(',')
			}

			builder.WriteString("\n  " + col.String())
		}

		builder.WriteString("\n)")
	}

	if domain.IsMemoryOptimized {
		builder.WriteString("\nWITH (MEMORY_OPTIMIZED = ON)")
	}

	builder.WriteString("\nGO")

	object.SetDefinition([]byte(builder.String()))

	return object, nil
}

const selectUserDefinedTypes = `
select userTypes.catalog, userTypes.type, userTypes.[schema], userTypes.parent_type, userTypes.max_length,
    userTypes.precision, userTypes.scale, userTypes.collation_name, userTypes.is_nullable,
    userTypes.is_table_type, userTypes.is_memory_optimized
from (
    select
        [catalog] = db_name(),
        [type] = types.name,
        [schema] = schema_name(types.schema_id),
        [parent_type] = st.name,
        [max_length] = iif(
            st.name in ('char', 'varchar', 'nchar', 'nvarchar', 'binary', 'varbinary'),
            iif(types.max_length = -1, 'max', cast(types.max_length as nvarchar(4))),
            iif(
                st.name = 'float',
                iif(
                    types.max_length != 53, cast(types.max_length as nvarchar(4)), null
                ),
                null
            )
        ),
        [precision] = iif(
            st.name in ('decimal', 'numeric'), types.precision, null
        ),
        [scale] = iif(
            st.name in ('decimal', 'numeric'), types.scale, null
        ),
        [collation_name] = types.collation_name,
        [is_nullable] = types.is_nullable,
        [is_table_type] = types.is_table_type,
        [is_memory_optimized] = cast(0 as bit)
    from sys.types as types
        inner join sys.types as st on (types.system_type_id = st.system_type_id)
            and st.system_type_id = st.user_type_id
    where types.is_user_defined != cast(0 as bit) and types.is_assembly_type = cast(0 as bit)
    union
    select
        [catalog] = db_name(),
        [type] = types.name,
        [schema] = schema_name(types.schema_id),
        [parent_type] = null,
        [max_length] = null,
        [precision] = null,
        [scale] = null,
        [collation_name] = types.collation_name,
        [is_nullable] = types.is_nullable,
        [is_table_type] = types.is_table_type,
        [is_memory_optimized] = types.is_memory_optimized
    from sys.table_types as types
        inner join sys.types as st on (types.system_type_id = st.system_type_id)
    where types.is_user_defined != cast(0 as bit) and types.is_assembly_type = cast(0 as bit)
) as userTypes
order by userTypes.[schema], userTypes.[type]
`
