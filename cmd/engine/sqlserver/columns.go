package sqlserver

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ColumnOption опция поля таблицы/табличного типа
type ColumnOption func(col Column)

// WithDefaultCollation collation по умолчанию
func WithDefaultCollation(collation string) ColumnOption {
	return func(col Column) {
		col.defaultCollation = collation
	}
}

// Column поле таблицы/табличного типа
type Column struct {
	ID                     int
	Name                   string
	description            sql.NullString
	TypeName               string
	TypeSchema             string
	IsUserDefinedType      bool
	maxLength              sql.NullString
	precision              sql.NullInt32
	scale                  sql.NullInt32
	collation              sql.NullString
	defaultCollation       string
	IsNullable             bool
	IsANSIPadded           bool
	IsRowGUIDCol           bool
	IsIdentity             bool
	identitySeedValue      sql.NullInt32
	identityIncrementValue sql.NullInt32
	isComputed             bool
	compute                sql.NullString
	IsFileStream           bool
	IsReplicated           bool
	IsNonSQLSubscribed     bool
	IsMergePublished       bool
	IsDTSReplicated        bool
	IsXMLDocument          bool
	def                    sql.NullString
	IsSparse               bool
	IsColumnSet            bool
	generateAlways         sql.NullString
	IsHidden               bool
	IsMasked               bool
}

// HasMaxLength проверяет, указана ли максимальная длина (в байтах) для типа
func (col Column) HasMaxLength() bool {
	return col.maxLength.Valid
}

// MaxLength максимальная длина значений типа данных поля
func (col Column) MaxLength() string {
	if col.maxLength.Valid {
		return col.maxLength.String
	}

	return ""
}

// HasPrecision проверяет, указана ли точность типа данных поля
func (col Column) HasPrecision() bool {
	return col.precision.Valid
}

// Precision точность типа данных поля
func (col Column) Precision() int {
	if col.precision.Valid {
		return int(col.precision.Int32)
	}

	return 0
}

// HasScale проверяет, указан ли масштаб для типа данных поля
func (col Column) HasScale() bool {
	return col.scale.Valid
}

// Scale возвращает масштаб для типа данных поля
func (col Column) Scale() int {
	if col.scale.Valid {
		return int(col.scale.Int32)
	}

	return 0
}

// HasCollation проверяет, указан ли collation
func (col Column) HasCollation() bool {
	return col.collation.Valid
}

// Collation
func (col Column) Collation() string {
	if col.collation.Valid {
		return col.collation.String
	}

	return ""
}

// IsComputed проверяет, является ли поле вычисляемым
func (col Column) IsComputed() bool {
	return col.isComputed && col.compute.Valid
}

// ComputeDefinition возвращает определение значения вычисляемого поля
func (col Column) ComputeDefinition() string {
	if col.IsComputed() {
		return col.compute.String
	}

	return ""
}

// HasDefaultDefinition проверяет, есть ли определение значения поля по умолчанию
func (col Column) HasDefaultDefinition() bool {
	return col.def.Valid
}

// DefaultDefinition возвращает определение значения поля по умолчанию
func (col Column) DefaultDefinition() string {
	if col.def.Valid {
		return col.def.String
	}

	return ""
}

// HasGenerateAlwaysDefinition проверяет, существует ли определение GENERATED ALWAYS... для поля периода
func (col Column) HasGenerateAlwaysDefinition() bool {
	return col.generateAlways.Valid
}

// GenerateAlwaysDefinition возвращает определение GENERATED ALWAYS... для поля периода в темпоральных таблицах
func (col Column) GenerateAlwaysDefinition() string {
	if col.generateAlways.Valid {
		return col.generateAlways.String
	}

	return ""
}

// HasDescription проверяет, есть ли описание поля
func (col Column) HasDescription() bool {
	return col.description.Valid
}

// Description возвращает описание поля
func (col Column) Description() string {
	if col.description.Valid {
		return col.description.String
	}

	return ""
}

// HasIdentitySeedValue проверяет, указано ли начальное значение для IDENTITY
func (col Column) HasIdentitySeedValue() bool {
	return col.identitySeedValue.Valid
}

// IdentitySeedValue возвращает начальное значение IDENTITY. Если не указано, то возвращает 1
func (col Column) IdentitySeedValue() int {
	if col.identitySeedValue.Valid {
		return int(col.identitySeedValue.Int32)
	}

	return 1
}

// HasIdentityIncrementValue проверяет, указано ли значение шага IDENTITY
func (col Column) HasIdentityIncrementValue() bool {
	return col.identityIncrementValue.Valid
}

// IdentityIncrementValue возвращает значение шага IDENTITY
func (col Column) IdentityIncrementValue() int {
	if col.identityIncrementValue.Valid {
		return int(col.identityIncrementValue.Int32)
	}

	return 1
}

// String возвращает определение поля
func (col Column) String() string {
	var builder strings.Builder

	builder.WriteString("[" + col.Name + "] ")

	if col.IsUserDefinedType {
		if strings.Trim(col.TypeSchema, " ") != "" {
			builder.WriteString("[" + col.TypeSchema + "].")
		}
	}

	builder.WriteString("[" + col.TypeName + "]")

	openedParentheses := col.HasMaxLength() || col.HasPrecision()

	if openedParentheses {
		builder.WriteRune('(')
	}

	if col.HasMaxLength() {
		builder.WriteString(col.MaxLength())
	} else {
		if col.HasPrecision() {
			precision := col.Precision()
			scale := col.Scale()

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

	if !col.IsNullable {
		builder.WriteString(" NOT NULL")
	}

	return builder.String()
}

// SetOptions устанавливает дополнительные опции поля таблицы/табличного типа
func (col Column) SetOptions(options ...ColumnOption) {
	for _, option := range options {
		option(col)
	}
}

// Columns поля
type Columns map[string]*Column

// SortedList возвращает срез полей
func (columns Columns) List() []*Column {
	length := len(columns)

	if length == 0 {
		return nil
	}

	list := make([]*Column, length)

	var index int

	for _, col := range columns {
		list[index] = col
		index++
	}

	return list
}

// ObjectColumns поля объектов (таблиц, табличных типов) БД
type ObjectColumns map[string]Columns

func (columns ObjectColumns) Append(objectSchema, objectName string, column *Column) error {
	name := SchemaAndObject(objectSchema, objectName, true)

	if strings.Trim(objectSchema, " ") == "" || strings.Trim(objectName, " ") == "" {
		return fmt.Errorf("impossible to identify the object %s", name)
	}

	if column == nil {
		return errors.New("column object cannot be nil")
	}

	if strings.Trim(column.Name, " ") == "" {
		return errors.New("impossible to identify the column")
	}

	if cols, ok := columns[name]; ok {
		cols[column.Name] = column
	} else {
		cols := make(Columns)
		cols[column.Name] = column

		columns[name] = cols
	}

	return nil
}

const selectColumns = `
select columns.catalog, columns.object_schema, columns.object_name, columns.column_id, columns.column_name,
    columns.column_description, columns.type_name, columns.type_schema, columns.is_user_defined_type,
    columns.max_length, columns.precision, columns.scale, columns.collation_name, columns.is_nullable,
    columns.is_ansi_padded, columns.is_rowguidcol, columns.is_identity,columns.seed_value, columns.increment_value,
    columns.is_computed, columns.computed_definition, columns.is_filestream, columns.is_replicated,
    columns.is_non_sql_subscribed, columns.is_merge_published, columns.is_dts_replicated, columns.is_xml_document,
    columns.default_object_definition, columns.is_sparse, columns.is_column_set, columns.generated_always,
    columns.is_hidden, columns.is_masked
from (
    select
        [catalog] = db_name(),
        [object_schema] = schema_name(types.schema_id),
        [object_name] = types.name,
        [object_id] = columns.object_id,
        [column_id] = columns.column_id,
        [column_name] = columns.name,
        [column_description] = cast(sep.value as nvarchar(2048)),
        [type_name] = st.name,
        [type_schema] = schema_name(st.schema_id),
        [is_user_defined_type] = st.is_user_defined,
        [max_length] = iif(
            st.name in ('char', 'varchar', 'nchar', 'nvarchar', 'binary', 'varbinary'),
            iif(columns.max_length = -1, 'max', cast(columns.max_length as nvarchar(4))),
            iif(
                st.name = 'float',
                iif(
                    columns.max_length != 53, cast(columns.max_length as nvarchar(4)), null
                ),
                null
            )
        ),
        [precision] = iif(
            st.name in ('decimal', 'numeric'), columns.precision, null
        ),
        [scale] = iif(
            st.name in ('decimal', 'numeric'), columns.scale, null
        ),
        [collation_name] = columns.collation_name,
        [is_nullable] = columns.is_nullable,
        [is_ansi_padded] = columns.is_ansi_padded,
        [is_rowguidcol] = columns.is_rowguidcol,
        [is_identity] = columns.is_identity,
        [seed_value] = ident.seed_value,
        [increment_value] = ident.increment_value,
        [is_computed] = columns.is_computed,
        [computed_definition] = cc.definition,
        [is_filestream] = columns.is_filestream,
        [is_replicated] = columns.is_replicated,
        [is_non_sql_subscribed] = columns.is_non_sql_subscribed,
        [is_merge_published] = columns.is_merge_published,
        [is_dts_replicated] = columns.is_dts_replicated,
        [is_xml_document] = columns.is_xml_document,
        [default_object_definition] = def.definition,
        [is_sparse] = columns.is_sparse,
        [is_column_set] = columns.is_column_set,
        [generated_always] = case columns.generated_always_type
            when 1 then 'GENERATED ALWAYS AS ROW START'
            when 2 then 'GENERATED ALWAYS AS ROW END'
            else null
        end,
        [is_hidden] = columns.is_hidden,
        [is_masked] = columns.is_masked
    from sys.table_types as types
        inner join sys.objects as objects on (types.type_table_object_id = objects.object_id)
            inner join sys.columns as columns on (objects.object_id = columns.object_id)
                inner join sys.types as st on (columns.user_type_id = st.user_type_id)
                left join sys.default_constraints as def on (columns.default_object_id = def.object_id)
                left join sys.computed_columns as cc on (columns.object_id = cc.object_id)
                    and (columns.column_id = cc.column_id)
                left join sys.identity_columns as ident on (columns.object_id = ident.object_id)
                    and (columns.column_id = ident.column_id)
                left join sys.extended_properties as sep on (columns.object_id = sep.major_id)
                    and (columns.column_id = sep.minor_id) and (sep.name = 'MS_Description')
                    and (sep.class = 1)
    where types.is_user_defined != cast(0 as bit) and types.is_assembly_type != cast(1 as bit)
    union
    select
        [catalog] = db_name(),
        [object_schema] = schema_name(tables.schema_id),
        [object_name] = tables.name,
        [object_id] = columns.object_id,
        [column_id] = columns.column_id,
        [column_name] = columns.name,
        [column_description] = cast(sep.value as nvarchar(2048)),
        [type_name] = st.name,
        [type_schema] = schema_name(st.schema_id),
        [is_user_defined_type] = st.is_user_defined,
        [max_length] = iif(
            st.name in ('char', 'varchar', 'nchar', 'nvarchar', 'binary', 'varbinary'),
            iif(columns.max_length = -1, 'max', cast(columns.max_length as nvarchar(4))),
            iif(
                st.name = 'float',
                iif(
                    columns.max_length != 53, cast(columns.max_length as nvarchar(4)), null
                ),
                null
            )
        ),
        [precision] = iif(
            st.name in ('decimal', 'numeric'), columns.precision, null
        ),
        [scale] = iif(
            st.name in ('decimal', 'numeric'), columns.scale, null
        ),
        [collation_name] = columns.collation_name,
        [is_nullable] = columns.is_nullable,
        [is_ansi_padded] = columns.is_ansi_padded,
        [is_rowguidcol] = columns.is_rowguidcol,
        [is_identity] = columns.is_identity,
        [seed_value] = ident.seed_value,
        [increment_value] = ident.increment_value,
        [is_computed] = columns.is_computed,
        [computed_definition] = cc.definition,
        [is_filestream] = columns.is_filestream,
        [is_replicated] = columns.is_replicated,
        [is_non_sql_subscribed] = columns.is_non_sql_subscribed,
        [is_merge_published] = columns.is_merge_published,
        [is_dts_replicated] = columns.is_dts_replicated,
        [is_xml_document] = columns.is_xml_document,
        [default_object_definition] = def.definition,
        [is_sparse] = columns.is_sparse,
        [is_column_set] = columns.is_column_set,
        [generated_always] = case columns.generated_always_type
            when 1 then 'GENERATED ALWAYS AS ROW START'
            when 2 then 'GENERATED ALWAYS AS ROW END'
            else null
        end,
        [is_hidden] = columns.is_hidden,
        [is_masked] = columns.is_masked
    from sys.tables as tables
        inner join sys.objects as objects on (tables.object_id = objects.object_id)
            inner join sys.columns as columns on (objects.object_id = columns.object_id)
                inner join sys.types as st on (columns.user_type_id = st.user_type_id)
                left join sys.default_constraints as def on (columns.default_object_id = def.object_id)
                left join sys.computed_columns as cc on (columns.object_id = cc.object_id)
                    and (columns.column_id = cc.column_id)
                left join sys.identity_columns as ident on (columns.object_id = ident.object_id)
                    and (columns.column_id = ident.column_id)
                left join sys.extended_properties as sep on (columns.object_id = sep.major_id)
                    and (columns.column_id = sep.minor_id) and (sep.name = 'MS_Description')
                    and (sep.class = 1)
    where tables.type = 'U'
) as columns
order by columns.catalog, columns.object_schema, columns.object_name, columns.column_id
`
