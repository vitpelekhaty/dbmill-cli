package sqlserver

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	str "github.com/vitpelekhaty/dbmill-cli/internal/pkg/strings"
)

// ColumnOption опция поля таблицы/табличного типа
type ColumnOption func(col *Column)

// WithDefaultCollation collation по умолчанию
func WithDefaultCollation(collation string) ColumnOption {
	return func(col *Column) {
		col.defaultCollation = collation
	}
}

// WithColumnOwner владелец поля (по умолчанию - OwnerTable)
func WithColumnOwner(owner ColumnOrIndexOwner) ColumnOption {
	return func(col *Column) {
		col.owner = owner
	}
}

// ColumnOrIndexOwner владелец поля/индекса
type ColumnOrIndexOwner byte

const (
	// OwnerTable владелец - таблица
	OwnerTable ColumnOrIndexOwner = iota
	// OwnerMemoryOptimizedTable владелец - memory-optimized таблица
	OwnerMemoryOptimizedTable
	// OwnerUserDefinedTableDataType владелец - пользовательский табличный тип
	OwnerUserDefinedTableDataType
	// OwnerAlterTable владелец - блок ALTER TABLE (для определения индексов)
	OwnerAlterTable
)

// Column поле таблицы/табличного типа
type Column struct {
	ID                            int
	Name                          string
	description                   sql.NullString
	TypeName                      string
	TypeSchema                    string
	IsUserDefinedType             bool
	maxLength                     sql.NullString
	precision                     sql.NullInt32
	scale                         sql.NullInt32
	collation                     sql.NullString
	IsNullable                    bool
	IsANSIPadded                  bool
	IsRowGUIDCol                  bool
	IsIdentity                    bool
	identitySeedValue             sql.NullInt32
	identityIncrementValue        sql.NullInt32
	isComputed                    bool
	isPersisted                   sql.NullBool
	compute                       sql.NullString
	IsFileStream                  bool
	IsReplicated                  bool
	IsNonSQLSubscribed            bool
	IsMergePublished              bool
	IsDTSReplicated               bool
	IsXMLDocument                 bool
	xmlSchemaCollectionSchemaName sql.NullString
	xmlSchemaCollectionName       sql.NullString
	defaultConstraint             sql.NullString
	defaultConstraintDefinition   sql.NullString
	IsSparse                      bool
	IsColumnSet                   bool
	generateAlways                sql.NullString
	IsHidden                      bool
	IsMasked                      bool
	maskingFunction               sql.NullString
	encryptionKey                 sql.NullString
	encryptionKeyDatabaseName     sql.NullString
	encryptionAlgorithm           sql.NullString
	encryptionType                sql.NullString

	defaultCollation string
	owner            ColumnOrIndexOwner
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

// HasDefaultConstraintDefinition проверяет, есть ли определение значения поля по умолчанию
func (col Column) HasDefaultConstraintDefinition() bool {
	return col.defaultConstraint.Valid && col.defaultConstraintDefinition.Valid
}

// DefaultConstraintDefinition возвращает определение значения поля по умолчанию
func (col Column) DefaultConstraintDefinition() string {
	if col.HasDefaultConstraintDefinition() {
		return col.defaultConstraintDefinition.String
	}

	return ""
}

// DefaultConstraintName возвращает имя default constraint
func (col Column) DefaultConstraintName() string {
	if col.HasDefaultConstraintDefinition() {
		return col.defaultConstraint.String
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

// HasMaskingFunction проверяет наличие masking function
func (col Column) HasMaskingFunction() bool {
	return col.IsMasked && col.maskingFunction.Valid
}

// MaskingFunction
func (col Column) MaskingFunction() string {
	if col.HasMaskingFunction() {
		return col.maskingFunction.String
	}

	return ""
}

// HasEncryption
func (col Column) HasEncryption() bool {
	return col.encryptionKey.Valid && col.encryptionType.Valid && col.encryptionAlgorithm.Valid
}

// EncryptionKey
func (col Column) EncryptionKey() string {
	if col.HasEncryption() {
		return col.encryptionKey.String
	}

	return ""
}

// EncryptionType
func (col Column) EncryptionType() string {
	if col.HasEncryption() {
		return col.encryptionType.String
	}

	return ""
}

// EncryptionAlgorithm
func (col Column) EncryptionAlgorithm() string {
	if col.HasEncryption() {
		return col.encryptionAlgorithm.String
	}

	return ""
}

// HasEncryptionDatabaseName
func (col Column) HasEncryptionDatabaseName() bool {
	return col.HasEncryption() && col.encryptionKeyDatabaseName.Valid
}

// EncryptionDatabaseName
func (col Column) EncryptionDatabaseName() string {
	if col.HasEncryptionDatabaseName() {
		return col.encryptionKeyDatabaseName.String
	}

	return ""
}

// HasXMLSchemaCollection
func (col Column) HasXMLSchemaCollection() bool {
	return strings.EqualFold(col.TypeName, "xml") && col.xmlSchemaCollectionName.Valid
}

// XMLSchemaCollection
func (col Column) XMLSchemaCollection() string {
	if col.HasXMLSchemaCollection() {
		var (
			schema string
			name   string
		)

		if col.xmlSchemaCollectionSchemaName.Valid {
			schema = col.xmlSchemaCollectionSchemaName.String
		}

		if col.xmlSchemaCollectionName.Valid {
			name = col.xmlSchemaCollectionName.String
		}

		return SchemaAndObject(schema, name, true)
	}

	return ""
}

// IsPersisted
func (col Column) IsPersisted() bool {
	if col.IsComputed() && col.isPersisted.Valid {
		return col.isPersisted.Valid
	}

	return false
}

// String возвращает определение поля
func (col Column) String() string {
	builder := str.NewBuilder("[" + col.Name + "]")

	dataTypeDefinition := col.dataTypeDefinition()

	if strings.Trim(dataTypeDefinition, " ") != "" {
		builder.WriteString(" " + dataTypeDefinition)
	}

	columnDefinition := col.columnDefinition()

	if strings.Trim(columnDefinition, " ") != "" {
		builder.WriteString(" " + columnDefinition)
	}

	return builder.String()
}

func (col Column) dataTypeDefinition() string {
	var builder str.Builder

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
				builder.WriteString(", " + strconv.Itoa(scale))
			}
		}
	}

	if openedParentheses {
		builder.WriteRune(')')
	}

	if col.HasXMLSchemaCollection() {
		builder.WriteSpace()

		if col.IsXMLDocument {
			builder.WriteString("DOCUMENT")
		} else {
			builder.WriteString("CONTENT")
		}

		builder.WriteSpace()
		builder.WriteString(col.XMLSchemaCollection())
	}

	return builder.String()
}

func (col Column) columnDefinition() string {
	if col.isComputed {
		return col.computedColumnDefinition()
	}

	switch col.owner {
	case OwnerTable:
		return col.columnDefinitionForTable()
	case OwnerMemoryOptimizedTable:
		return col.columnDefinitionForMemoryOptimizedTable()
	case OwnerUserDefinedTableDataType:
		return col.columnDefinitionForUserDefinedTableType()
	default:
		return ""
	}
}

func (col Column) computedColumnDefinition() string {
	return ""
}

func (col Column) columnDefinitionForTable() string {
	var builder str.Builder

	if col.IsFileStream {
		builder.WriteString("FILESTREAM")
	}

	if col.HasCollation() {
		collation := col.Collation()

		if collation != col.defaultCollation {
			builder.WriteSpace()
			builder.WriteString("COLLATE " + collation)
		}
	}

	if col.IsSparse {
		builder.WriteSpace()
		builder.WriteString("SPARSE")
	}

	if col.HasMaskingFunction() {
		maskedColumnOption := fmt.Sprintf(`MASKED WITH (FUNCTION = '%s')`, col.MaskingFunction())

		builder.WriteSpace()
		builder.WriteString(maskedColumnOption)
	}

	if col.HasDefaultConstraintDefinition() {
		defaultDefinition := fmt.Sprintf(`CONSTRAINT [%s] DEFAULT %s`, col.DefaultConstraintName(),
			col.DefaultConstraintDefinition())

		builder.WriteSpace()
		builder.WriteString(defaultDefinition)
	}

	if col.IsIdentity {
		builder.WriteSpace()
		builder.WriteString("IDENTITY")

		seed := col.IdentitySeedValue()
		increment := col.IdentityIncrementValue()

		if seed > 1 || increment > 1 {
			identityValues := fmt.Sprintf("(%d, %d)", seed, increment)
			builder.WriteString(identityValues)
		}

		if !col.IsReplicated {
			builder.WriteSpace()
			builder.WriteString("NOT FOR REPLICATION")
		}
	}

	if col.HasGenerateAlwaysDefinition() {
		builder.WriteSpace()
		builder.WriteString(col.GenerateAlwaysDefinition())

		if col.IsHidden {
			builder.WriteSpace()
			builder.WriteString("HIDDEN")
		}
	}

	if !col.IsNullable && !col.IsIdentity {
		builder.WriteSpace()
		builder.WriteString("NOT NULL")
	}

	if col.IsRowGUIDCol {
		builder.WriteSpace()
		builder.WriteString("ROWGUIDCOL")
	}

	if col.HasEncryption() {
		builder.WriteSpace()

		encryptionOptions := fmt.Sprintf(`COLUMN_ENCRYPTION_KEY = %s, ENCRYPTION_TYPE = %s, ALGORITHM = '%s'`,
			col.EncryptionKey(), col.EncryptionType(), col.EncryptionAlgorithm())
		builder.WriteString("ENCRYPTED WITH (" + encryptionOptions + ")")
	}

	return builder.String()
}

func (col Column) columnDefinitionForMemoryOptimizedTable() string {
	// TODO: требуется тестирование описаний полей для memory-optimized таблиц
	var builder str.Builder

	if col.HasCollation() {
		collation := col.Collation()

		if collation != col.defaultCollation {
			builder.WriteSpace()
			builder.WriteString("COLLATE " + collation)
		}
	}

	if col.HasGenerateAlwaysDefinition() {
		builder.WriteSpace()
		builder.WriteString(col.GenerateAlwaysDefinition())

		if col.IsHidden {
			builder.WriteSpace()
			builder.WriteString("HIDDEN")
		}
	}

	if !col.IsNullable && !col.IsIdentity {
		builder.WriteSpace()
		builder.WriteString("NOT NULL")
	}

	if col.HasDefaultConstraintDefinition() {
		defaultDefinition := fmt.Sprintf(`CONSTRAINT [%s] DEFAULT %s`, col.DefaultConstraintName(),
			col.DefaultConstraintDefinition())

		builder.WriteSpace()
		builder.WriteString(defaultDefinition)
	}

	if col.IsIdentity {
		builder.WriteSpace()
		builder.WriteString("IDENTITY")

		seed := col.IdentitySeedValue()
		increment := col.IdentityIncrementValue()

		if seed > 1 || increment > 1 {
			identityValues := fmt.Sprintf("(%d, %d)", seed, increment)
			builder.WriteString(identityValues)
		}

		if !col.IsReplicated {
			builder.WriteSpace()
			builder.WriteString("NOT FOR REPLICATION")
		}
	}

	return builder.String()
}

func (col Column) columnDefinitionForUserDefinedTableType() string {
	var builder str.Builder

	if col.HasCollation() {
		collation := col.Collation()

		if collation != col.defaultCollation {
			builder.WriteSpace()
			builder.WriteString("COLLATE " + collation)
		}
	}

	if !col.IsNullable && !col.IsIdentity {
		builder.WriteSpace()
		builder.WriteString("NOT NULL")
	}

	if col.HasDefaultConstraintDefinition() {
		defaultDefinition := fmt.Sprintf(`CONSTRAINT [%s] DEFAULT %s`, col.DefaultConstraintName(),
			col.DefaultConstraintDefinition())

		builder.WriteSpace()
		builder.WriteString(defaultDefinition)
	}

	if col.IsIdentity {
		builder.WriteSpace()
		builder.WriteString("IDENTITY")

		seed := col.IdentitySeedValue()
		increment := col.IdentityIncrementValue()

		if seed > 1 || increment > 1 {
			identityValues := fmt.Sprintf("(%d, %d)", seed, increment)
			builder.WriteString(identityValues)
		}
	}

	if col.IsRowGUIDCol {
		builder.WriteSpace()
		builder.WriteString("ROWGUIDCOL")
	}

	return builder.String()
}

// SetOptions устанавливает дополнительные опции поля таблицы/табличного типа
func (col *Column) SetOptions(options ...ColumnOption) {
	for _, option := range options {
		option(col)
	}
}

// Columns поля
type Columns map[string]*Column

// Slice возвращает срез полей
func (columns Columns) Slice() []*Column {
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
with obj (schema_id, object_id, object_name) as (
    select ttypes.schema_id, ttypes.type_table_object_id as object_id, ttypes.name as object_name
    from sys.table_types as ttypes
    where ttypes.is_user_defined != cast(0 as bit) and ttypes.is_assembly_type != cast(1 as bit)
    union
    select tables.schema_id, tables.object_id, tables.name as object_name
    from sys.tables as tables
    where tables.type = 'U'
)
select columns.catalog, columns.object_schema, columns.object_name, columns.column_id, columns.column_name,
    columns.column_description, columns.type_name, columns.type_schema, columns.is_user_defined_type,
    columns.max_length, columns.precision, columns.scale, columns.collation_name, columns.is_nullable,
    columns.is_ansi_padded, columns.is_rowguidcol, columns.is_identity,columns.seed_value, columns.increment_value,
    columns.is_computed, columns.is_persisted, columns.computed_definition, columns.is_filestream, columns.is_replicated,
    columns.is_non_sql_subscribed, columns.is_merge_published, columns.is_dts_replicated, columns.is_xml_document,
    columns.xml_schema_collection_schema, columns.xml_schema_collection_name, columns.default_constraint,
    columns.default_object_definition, columns.is_sparse, columns.is_column_set, columns.generated_always,
    columns.is_hidden, columns.is_masked, columns.masking_function, columns.encryption_key, columns.encryption_type,
    columns.encryption_algorithm, columns.encryption_key_database_name
from (
    select
        [catalog] = db_name(),
        [object_schema] = schema_name(o.schema_id),
        [object_name] = o.object_name,
        [object_id] = columns.object_id,
        [column_id] = columns.column_id,
        [column_name] = columns.name,
        [column_description] = cast(sep.value as nvarchar(2048)),
        [type_name] = st.name,
        [type_schema] = schema_name(st.schema_id),
        [is_user_defined_type] = st.is_user_defined,
        [max_length] = iif(
            st.name in ('char', 'varchar', 'nchar', 'nvarchar', 'binary', 'varbinary'),
            iif(
                    columns.max_length = -1,
                    'max',
                    iif(
                            st.name in ('nchar', 'nvarchar'),
                            cast(columns.max_length / 2 as nvarchar(4)),
                            cast(columns.max_length as nvarchar(4))
                        )
            ),
            iif(
                st.name = 'float',
                iif(columns.max_length != 53, cast(columns.max_length as nvarchar(4)), null),
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
        [is_persisted] = cc.is_persisted,
        [computed_definition] = cc.definition,
        [is_filestream] = columns.is_filestream,
        [is_replicated] = columns.is_replicated,
        [is_non_sql_subscribed] = columns.is_non_sql_subscribed,
        [is_merge_published] = columns.is_merge_published,
        [is_dts_replicated] = columns.is_dts_replicated,
        [is_xml_document] = columns.is_xml_document,
        [xml_schema_collection_schema] = schema_name(xsc.schema_id),
        [xml_schema_collection_name] = xsc.name,
        [default_constraint] = def.name,
        [default_object_definition] = def.definition,
        [is_sparse] = columns.is_sparse,
        [is_column_set] = columns.is_column_set,
        [generated_always] = case columns.generated_always_type
            when 1 then 'GENERATED ALWAYS AS ROW START'
            when 2 then 'GENERATED ALWAYS AS ROW END'
            else null
        end,
        [is_hidden] = columns.is_hidden,
        [is_masked] = columns.is_masked,
        [masking_function] = mc.masking_function,
        [encryption_key] = cek.name,
        [encryption_type] = case columns.encryption_type
            when 1 then N'DETERMINISTIC'
            when 2 then N'RANDOMIZED'
            else ''
        end,
        [encryption_algorithm] = columns.encryption_algorithm_name,
        [encryption_key_database_name] = columns.column_encryption_key_database_name

    from obj as o
        inner join sys.objects as objects on (o.object_id = objects.object_id)
            inner join sys.columns as columns on (objects.object_id = columns.object_id)
                inner join sys.types as st on (columns.user_type_id = st.user_type_id)
                left join sys.default_constraints as def on (columns.default_object_id = def.object_id)
                left join sys.computed_columns as cc on (columns.object_id = cc.object_id)
                    and (columns.column_id = cc.column_id)
                left join sys.identity_columns as ident on (columns.object_id = ident.object_id)
                    and (columns.column_id = ident.column_id)
                left join sys.masked_columns as mc on (columns.object_id = mc.object_id)
                    and (columns.column_id = mc.column_id) and (mc.is_masked != cast(0 as bit))
                left join sys.extended_properties as sep on (columns.object_id = sep.major_id)
                    and (columns.column_id = sep.minor_id) and (sep.name = 'MS_Description')
                    and (sep.class = 1)
                left join sys.column_encryption_keys as cek on (columns.column_encryption_key_id = cek.column_encryption_key_id)
                left join sys.xml_schema_collections as xsc on (columns.xml_collection_id = xsc.xml_collection_id)
) as columns
order by columns.catalog, columns.object_schema, columns.object_name, columns.column_id
`
