package sqlserver

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	str "github.com/vitpelekhaty/dbmill-cli/internal/pkg/strings"
)

// IndexType тип индекса
type IndexType byte

const (
	// IndexTypeCustomIndex пользовательский индекс
	IndexTypeCustomIndex IndexType = iota
	// IndexTypePrimaryKey первичный ключ
	IndexTypePrimaryKey
	// IndexTypeUnique уникальный индекс
	IndexTypeUnique
)

// ColumnOption опция поля таблицы/табличного типа
type IndexOption func(col *Index)

// WithIndexOwner владелец индекса (по умолчанию - OwnerTable)
func WithIndexOwner(owner ColumnOrIndexOwner) IndexOption {
	return func(col *Index) {
		col.owner = owner
	}
}

// Index определение индекса
// (https://docs.microsoft.com/en-us/sql/relational-databases/system-catalog-views/sys-indexes-transact-sql)
type Index struct {
	// Name наименование
	Name string
	// Type тип индекса (HEAP, CLUSTERED, NONCLUSTERED etc)
	Type string
	// IsUnique уникальный индекс
	IsUnique bool
	// IsPrimaryKey первичный ключ
	IsPrimaryKey bool
	// IsUniqueConstraint индекс является частью ограничения UNIQUE
	IsUniqueConstraint bool
	// IgnoreDupKey значение параметра IGNORE_DUP_KEY
	IgnoreDupKey bool
	// FillFactor значение FILLFACTOR
	FillFactor byte
	// IsPadded значение параметра PADINDEX
	IsPadded bool
	// IsDisabled индекс отключен
	IsDisabled bool
	// IsHypothetical индекс является гипотетическим
	IsHypothetical bool
	// IsIgnoredInOptimization
	IsIgnoredInOptimization bool
	// AllowRowLocks индекс допускает блокировку строк
	AllowRowLocks bool
	// AllowPageLocks индекс допускает блокировку страниц
	AllowPageLocks bool
	// SuppressDupKeyMessages
	SuppressDupKeyMessages bool
	// AutoCreated индекс создан автоматической настройкой
	AutoCreated bool
	// OptimizeForSequentialKey для индекса включена оптимизация вставки последней страницы
	OptimizeForSequentialKey bool
	// Columns ключевые поля в индексе
	Columns IndexedColumns
	// IncludedColumns неключевые поля, включенные в индекс
	IncludedColumns IndexedColumns

	hasFilter        bool
	filterDefinition sql.NullString
	bucketCount      sql.NullInt64
	description      sql.NullString

	owner ColumnOrIndexOwner
}

// BucketCount возвращает число контейнеров, которые необходимо создать в хэш-индексе
func (index Index) BucketCount() int {
	if index.bucketCount.Valid {
		return int(index.bucketCount.Int64)
	}

	return 0
}

// HasFilter индекс с фильтром
func (index Index) HasFilter() bool {
	return index.hasFilter
}

// FilterDefinition возвращает выражение для подмножества строк, включенного в фильтруемый индекс
func (index Index) FilterDefinition() string {
	if index.hasFilter && index.filterDefinition.Valid {
		return index.filterDefinition.String
	}

	return ""
}

// Description описание индекса
func (index Index) Description() string {
	if index.description.Valid {
		return index.description.String
	}

	return ""
}

// HasDescription проверяет, есть ли у индекса описание
func (index Index) HasDescription() bool {
	return index.description.Valid
}

// IndexType тип индекса
func (index Index) IndexType() IndexType {
	if index.IsPrimaryKey {
		return IndexTypePrimaryKey
	}

	if index.IsUnique {
		return IndexTypeUnique
	}

	return IndexTypeCustomIndex
}

// IsClustered кластерный индекс
func (index Index) IsClustered() bool {
	return index.IsType("CLUSTERED")
}

// IsNonClustered некластерный индекс
func (index Index) IsNonClustered() bool {
	return index.IsType("NONCLUSTERED")
}

// IsHash хеш-индекс
func (index Index) IsHash() bool {
	return index.IsType("HASH")
}

// IsType проверяет, является ли индекс индексом указанного типа
func (index Index) IsType(typ string) bool {
	return strings.EqualFold(index.Type, typ)
}

// Flags флаги создания индекса
func (index Index) Flags() []string {
	flags := make([]string, 0)

	if index.IsHash() {
		flags = append(flags, fmt.Sprintf("BUCKET_COUNT = %d", index.BucketCount()))
	}

	if index.IgnoreDupKey {
		flags = append(flags, "IGNORE_DUP_KEY = ON")
	}

	return flags
}

// SetOptions устанавливает дополнительные опции поля таблицы/табличного типа
func (index *Index) SetOptions(options ...IndexOption) {
	for _, option := range options {
		option(index)
	}
}

// String возвращает определение индекса
func (index Index) String() string {
	switch index.owner {
	case OwnerUserDefinedTableDataType:
		return index.constraintDefinitionForTableType()
	case OwnerMemoryOptimizedTable:
		return index.constraintDefinitionForMemoryOptimizedTable()
	case OwnerAlterTable:
		return index.constraintDefinitionForAlterTableBlock()
	default:
		return index.constraintDefinitionForTable()
	}
}

func (index Index) constraintDefinitionForTable() string {
	var builder str.Builder
	return builder.String()
}

func (index Index) constraintDefinitionForTableType() string {
	var builder str.Builder

	if index.IndexType() == IndexTypePrimaryKey {
		builder.WriteString("PRIMARY KEY")
	} else {
		builder.WriteString("INDEX")
	}

	builder.WriteSpace()
	builder.WriteString("[" + index.Name + "]")

	if index.IndexType() == IndexTypeUnique {
		builder.WriteSpace()
		builder.WriteString("UNIQUE")
	}

	if !index.IsNonClustered() {
		builder.WriteSpace()
		builder.WriteString(index.Type)
	}

	columns := index.Columns.Slice()

	if len(columns) > 0 {
		sort.Slice(columns, func(i, j int) bool {
			return columns[i].KeyOrdinal < columns[j].KeyOrdinal
		})

		ct := columns.Join(true, ", ")

		builder.WriteSpace()
		builder.WriteString("(" + ct + ")")
	}

	flags := index.Flags()

	if len(flags) > 0 {
		sort.Strings(flags)

		builder.WriteSpace()
		builder.WriteString(fmt.Sprintf("WITH (%s)", strings.Join(flags, ", ")))
	}

	includedColumns := index.IncludedColumns.Slice()

	if len(includedColumns) > 0 {
		sort.Slice(includedColumns, func(i, j int) bool {
			return includedColumns[i].KeyOrdinal < includedColumns[j].KeyOrdinal
		})

		ict := includedColumns.Join(true, ", ")

		builder.WriteSpace()
		builder.WriteString("INCLUDE (" + ict + ")")
	}

	return builder.String()
}

func (index Index) constraintDefinitionForMemoryOptimizedTable() string {
	var builder str.Builder
	return builder.String()
}

func (index Index) constraintDefinitionForAlterTableBlock() string {
	var builder str.Builder
	return builder.String()
}

// IndexedColumn индексируемые поля
// (https://docs.microsoft.com/en-us/sql/relational-databases/system-catalog-views/sys-index-columns-transact-sql)
type IndexedColumn struct {
	// ID поля в индексе
	ID int
	// Name наименование поля
	Name string
	// IsDescendingKey насправление сортировки - по убыванию
	IsDescendingKey bool
	// KeyOrdinal порядковый номер (нумерация начинается с 1) внутри набора ключевых столбцов
	KeyOrdinal int
	// PartitionOrdinal порядковый номер (нумерация начинается с 1) внутри набора столбцов секционирования
	PartitionOrdinal int
	// ColumnStoreOrderOrdinal порядковый номер (от 1) в наборе столбцов в упорядоченном кластеризованном
	// индексе columnstore
	ColumnStoreOrderOrdinal int
}

// IndexedColumnsSlice срез индексируемых полей
type IndexedColumnsSlice []*IndexedColumn

func (columns IndexedColumnsSlice) Join(useBrackets bool, sep string) string {
	if len(columns) == 0 {
		return ""
	}

	col := make([]string, len(columns))
	var columnName string

	for index, column := range columns {
		columnName = column.Name

		if useBrackets {
			columnName = "[" + columnName + "]"
		}

		if column.IsDescendingKey {
			col[index] = columnName + " DESC"
		} else {
			col[index] = columnName
		}
	}

	return strings.Join(col, sep)
}

// IndexedColumns тип справочника индексируемых полей. Ключ справочника - наименование поля
type IndexedColumns map[string]*IndexedColumn

// Slice возвращает срез индексируемых полей
func (columns IndexedColumns) Slice() IndexedColumnsSlice {
	if len(columns) == 0 {
		return nil
	}

	out := make([]*IndexedColumn, len(columns))
	var i int

	for _, column := range columns {
		out[i] = column
		i++
	}

	return out
}

// Indexes тип справочника индексов объекта БД. Ключ справочника - наименование индекса
type Indexes map[string]*Index

// PrimaryKeys справочник первичных ключей объекта БД
func (indexes Indexes) PrimaryKeys() Indexes {
	return indexes.filterByIndexType(IndexTypePrimaryKey)
}

// UniqueIndexes справочник уникальных индексов объекта БД
func (indexes Indexes) UniqueIndexes() Indexes {
	return indexes.filterByIndexType(IndexTypeUnique)
}

// CustomIndexes справочник пользовательских индексов объекта БД
func (indexes Indexes) CustomIndexes() Indexes {
	return indexes.filterByIndexType(IndexTypeCustomIndex)
}

// filterByIndexType возвращает выбранные индексы указанного типа
func (indexes Indexes) filterByIndexType(indexType IndexType) Indexes {
	out := make(Indexes)

	for indexName, index := range indexes {
		if index.IndexType() == indexType {
			out[indexName] = index
		}
	}

	return out
}

// Slice возвращает срез индексов объекта БД
func (indexes Indexes) Slice() []*Index {
	if len(indexes) == 0 {
		return nil
	}

	out := make([]*Index, len(indexes))
	var i int

	for _, index := range indexes {
		out[i] = index
		i++
	}

	return out
}

// ObjectsIndexes тип справочника определений индексов объектов БД. Ключ справочника - наименование объекта БД
type ObjectsIndexes map[string]Indexes

// ForeignKey определение внешнего ключа
type ForeignKey struct {
	// Name наименование внешнего ключа
	Name string
	// ReferencedObjectSchema схема объекта БД, на который ссылается ключ
	ReferencedObjectSchema string
	// ReferencedObjectName наименование объекта БД, на который ссылается ключ
	ReferencedObjectName string
	// IsDisabled ключ отключен
	IsDisabled bool
	// IsNotForReplication ключ создан с опцией NOT FOR REPLICATION
	IsNotForReplication bool
	// IsNotTrusted ограничение ключа не проверено системой
	IsNotTrusted bool
	// DeleteReferentialAction действие при удалении значения, на которое ссылается ключ
	DeleteReferentialAction string
	// UpdateReferentialAction действие при обновлении значения, на которое ссылается ключ
	UpdateReferentialAction string
	// ColumnsReferences ссылки ключевых полей
	ColumnsReferences map[string]*ColumnReference

	// description описание внешнего ключа
	description sql.NullString
}

// HasDescription проверяет наличие описания о внешнего ключа
func (fk ForeignKey) HasDescription() bool {
	return fk.description.Valid
}

// Description описание внешнего ключа
func (fk ForeignKey) Description() string {
	if fk.description.Valid {
		return fk.description.String
	}

	return ""
}

// ColumnReference ссылка поля на поле внешнего объекта
type ColumnReference struct {
	// ID идентификатор ссылки
	ID int
	// Column поле объекта БД
	Column string
	// ReferencedColumn поле внешнего объекта БД
	ReferencedColumn string
}

// ForeignKeys тип справочника определений внешних ключей. Ключ справочника - наименование внешнего ключа
type ForeignKeys map[string]*ForeignKey

// Slice возвращает срез внешних ключей
func (keys ForeignKeys) Slice() []*ForeignKey {
	if len(keys) == 0 {
		return nil
	}

	out := make([]*ForeignKey, len(keys))
	var i int

	for _, index := range keys {
		out[i] = index
		i++
	}

	return out
}

// ObjectsForeignKeys тип справочника определений внешних ключей объектов. Ключ справочника - наименование объекта БД
type ObjectsForeignKeys map[string]ForeignKeys

const selectIndexes = `
select indexes.catalog, indexes.[schema], indexes.object_name, indexes.index_name, indexes.index_type,
    indexes.is_unique, indexes.is_primary_key, indexes.is_unique_constraint, indexes.ignore_dup_key,
    indexes.fill_factor, indexes.is_padded, indexes.is_disabled, indexes.is_hypothetical,
    indexes.is_ignored_in_optimization, indexes.allow_row_locks, indexes.allow_page_locks,
    indexes.suppress_dup_key_messages, indexes.auto_created, indexes.optimize_for_sequential_key, indexes.has_filter,
    indexes.filter_definition, indexes.index_column_id, indexes.column_name, indexes.is_descending_key,
    indexes.is_included_column, indexes.key_ordinal, indexes.partition_ordinal, indexes.column_store_order_ordinal,
    indexes.bucket_count, indexes.description
from (
    select
        [catalog] = db_name(),
        [schema] = iif(objects.type = 'TT', schema_name(table_types.schema_id), schema_name(objects.schema_id)),
        [object_name] = iif(objects.type = 'TT', table_types.name, objects.name),
        [object_type] = objects.type + ' - ' + objects.type_desc,
        [index_id] = indexes.index_id,
        [index_name] = indexes.name,
        [index_type] = indexes.type_desc,
        [is_unique] = indexes.is_unique,
        [is_primary_key] = indexes.is_primary_key,
        [is_unique_constraint] = indexes.is_unique_constraint,
        [ignore_dup_key] = indexes.ignore_dup_key,
        [fill_factor] = indexes.fill_factor,
        [is_padded] = indexes.is_padded,
        [is_disabled] = indexes.is_disabled,
        [is_hypothetical] = indexes.is_hypothetical,
        [is_ignored_in_optimization] = indexes.is_ignored_in_optimization,
        [allow_row_locks] = indexes.allow_row_locks,
        [allow_page_locks] = indexes.allow_page_locks,
        [suppress_dup_key_messages] = indexes.suppress_dup_key_messages,
        [auto_created] = indexes.auto_created,
        [optimize_for_sequential_key] = indexes.optimize_for_sequential_key,
        [has_filter] = indexes.has_filter,
        [filter_definition] = indexes.filter_definition,
        [index_column_id] = index_columns.index_column_id,
        [column_name] = columns.name,
        [is_descending_key] = index_columns.is_descending_key,
        [is_included_column] = index_columns.is_included_column,
        [key_ordinal] = index_columns.key_ordinal,
        [partition_ordinal] = index_columns.partition_ordinal,
        [column_store_order_ordinal] = index_columns.column_store_order_ordinal,
        [bucket_count] = hash_indexes.bucket_count,
        [description] = cast(prop.value as nvarchar(2048))

    from sys.indexes as indexes
        inner join sys.objects as objects on (indexes.object_id = objects.object_id)
            and (objects.type in ('U', 'V', 'TF', 'TT'))
            left join sys.table_types as table_types on (objects.object_id = table_types.type_table_object_id)
            left join sys.extended_properties as prop on (objects.object_id = prop.major_id)
                and (prop.minor_id = indexes.index_id) and (prop.name = 'MS_Description') and (prop.class = 7)
        inner join sys.index_columns as index_columns on (indexes.object_id = index_columns.object_id)
            and (indexes.index_id = index_columns.index_id)
            inner join sys.columns as columns on (index_columns.object_id = columns.object_id)
                and (index_columns.column_id = columns.column_id)
        left join sys.hash_indexes as hash_indexes on (indexes.object_id = hash_indexes.object_id)
            and (indexes.index_id = hash_indexes.index_id)
) as indexes
order by indexes.catalog, indexes.[schema], indexes.object_name, indexes.index_id, indexes.index_column_id`

const selectForeignKeys = `
select fk.catalog, fk.foreign_key_name, fk.constraint_column_id, fk.parent_object_schema, fk.parent_object_name,
    fk.parent_column_name, fk.referenced_object_schema, fk.referenced_object_name, fk.referenced_columns_name,
    fk.is_disabled, fk.is_not_for_replication, fk.is_not_trusted, fk.delete_referential_action,
    fk.update_referential_action, fk.description
from (
    select
        [catalog] = db_name(),
        [foreign_key_name] = constraint_objects.name,
        [constraint_column_id] = fk_columns.constraint_column_id,
        [parent_object_schema] = schema_name(parent_objects.schema_id),
        [parent_object_name] = parent_objects.name,
        [parent_column_name] = parent_columns.name,
        [referenced_object_schema] = schema_name(referenced_objects.schema_id),
        [referenced_object_name] = referenced_objects.name,
        [referenced_columns_name] = referenced_columns.name,
        [key_index_id] = fk.key_index_id,
        [is_disabled] = fk.is_disabled,
        [is_not_for_replication] = fk.is_not_for_replication,
        [is_not_trusted] = fk.is_not_trusted,
        [delete_referential_action] = case fk.delete_referential_action
            when 1 then 'CASCADE'
            when 2 then 'SET NULL'
            when 3 then 'SET DEFAULT'
            else 'NO ACTION'
        end,
        [update_referential_action] = case fk.update_referential_action
            when 1 then 'CASCADE'
            when 2 then 'SET NULL'
            when 3 then 'SET DEFAULT'
            else 'NO ACTION'
        end,
        [description] = cast(prop.value as nvarchar(2048))

    from sys.foreign_key_columns as fk_columns
        inner join sys.objects as constraint_objects on (fk_columns.constraint_object_id = constraint_objects.object_id)
            inner join sys.foreign_keys as fk on (constraint_objects.object_id = fk.object_id)
                left join sys.extended_properties as prop on (fk.object_id = prop.major_id)
                    and (prop.minor_id = fk.key_index_id) and (prop.name = 'MS_Description') and (prop.class = 7)
        inner join sys.columns as parent_columns on (fk_columns.parent_object_id = parent_columns.object_id)
            and (fk_columns.parent_column_id = parent_columns.column_id)
            inner join sys.objects as parent_objects on (parent_columns.object_id = parent_objects.object_id)
        inner join sys.columns as referenced_columns on (fk_columns.referenced_object_id = referenced_columns.object_id)
            and (fk_columns.referenced_column_id = referenced_columns.column_id)
            inner join sys.objects as referenced_objects on (referenced_columns.object_id = referenced_objects.object_id)
) as fk
order by fk.catalog, fk.parent_object_schema, fk.parent_object_name, fk.foreign_key_name, fk.constraint_column_id
`
