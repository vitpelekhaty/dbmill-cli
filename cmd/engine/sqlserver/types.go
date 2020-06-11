package sqlserver

import (
	"context"
	"database/sql"
)

type userDefinedType struct {
	catalog           string
	typeName          string
	schema            string
	parentTypeName    string
	maxLength         sql.NullString
	precision         sql.NullInt32
	scale             sql.NullInt32
	collation         sql.NullString
	isNullable        bool
	isTableType       bool
	isMemoryOptimized bool
}

// SchemaAndObject возвращает наименование объекта в формате %schema%.%name%
func (udt userDefinedType) SchemaAndName(useBrackets bool) string {
	return SchemaAndObject(udt.schema, udt.typeName, useBrackets)
}

// HasMaxLength проверяет, указана ли максимальная длина (в байтах) для типа
func (udt userDefinedType) HasMaxLength() bool {
	return udt.maxLength.Valid
}

// MaxLength максимальная длина (в байтах) типа
func (udt userDefinedType) MaxLength() string {
	if udt.maxLength.Valid {
		return udt.maxLength.String
	}

	return ""
}

// HasPrecision проверяет, указана ли точность
func (udt userDefinedType) HasPrecision() bool {
	return udt.precision.Valid
}

// Precision точность
func (udt userDefinedType) Precision() int {
	if udt.precision.Valid {
		return int(udt.precision.Int32)
	}

	return 0
}

// HasScale проверяет, указан ли масштаб
func (udt userDefinedType) HasScale() bool {
	return udt.scale.Valid
}

// Scale масштаб
func (udt userDefinedType) Scale() int {
	if udt.scale.Valid {
		return int(udt.scale.Int32)
	}

	return 0
}

// HasCollation провеяет, указан ли collation
func (udt userDefinedType) HasCollation() bool {
	return udt.collation.Valid
}

// Collation
func (udt userDefinedType) Collation() string {
	if udt.collation.Valid {
		return udt.collation.String
	}

	return ""
}

type userDefinedTypes map[string]*userDefinedType

func (types userDefinedTypes) append(udt *userDefinedType) {
	if udt == nil {
		return
	}

	types[udt.SchemaAndName(true)] = udt
}

func (command *ScriptsFolderCommand) userTypes() (userDefinedTypes, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	stmt, err := command.engine.db.PrepareContext(ctx, selectUserDefinedTypes)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		catalog           string
		typeName          string
		schema            string
		parentTypeName    string
		maxLength         sql.NullString
		precision         sql.NullInt32
		scale             sql.NullInt32
		collation         sql.NullString
		isNullable        bool
		isTableType       bool
		isMemoryOptimized bool
	)

	types := make(userDefinedTypes)

	for rows.Next() {
		err = rows.Scan(&catalog, &typeName, &schema, &parentTypeName, &maxLength, &precision, &scale, &collation,
			&isNullable, &isTableType, &isMemoryOptimized)

		if err != nil {
			return nil, err
		}

		t := &userDefinedType{
			catalog:           catalog,
			typeName:          typeName,
			schema:            schema,
			parentTypeName:    parentTypeName,
			maxLength:         maxLength,
			precision:         precision,
			scale:             scale,
			collation:         collation,
			isNullable:        isNullable,
			isTableType:       isTableType,
			isMemoryOptimized: isMemoryOptimized,
		}

		types.append(t)
	}

	return types, nil
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
        [parent_type] = st.name,
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
