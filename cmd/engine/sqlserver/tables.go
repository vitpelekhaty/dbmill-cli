package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

// DataSpace пространство данных
type DataSpace struct {
	// Name имя пространства данных
	Name string
	// Type тип пространства данных
	// (FILESTREAM_DATA_FILEGROUP | MEMORY_OPTIMIZED_DATA_FILEGROUP | PARTITION_SCHEME | ROWS_FILEGROUP)
	Type string
	// IsDefault пространство данных по умолчанию
	IsDefault bool
}

// Table таблица
type Table struct {
	// Catalog название базы данных
	Catalog string
	// Schema схема таблицы
	Schema string
	// Name наименование таблицы
	Name string
	// LOBDataSpace пространство имен больших двоичных объектов
	LOBDataSpace *DataSpace
	// LockOnBulkLoad блокировка при массовом обновлении (если отключена (по умолчанию), то при массовой загрузке
	// данных осуществляется блокировка строк)
	LockOnBulkLoad bool
	// UsesANSINulls таблица создана с параметром ANSI_NULLS = ON
	UsesANSINulls bool
	// IsReplicated таблица создана при репликации транзакции или снимка
	IsReplicated bool
	// HasReplicationFilter для таблицы имеется фильтр репликации
	HasReplicationFilter bool
	// IsMergePublished таблица создана путем репликации слиянием
	IsMergePublished bool
	// IsSyncTranSubscribed таблица используется в немедленно обновляемой подписке
	IsSyncTranSubscribed bool
	// HasUncheckedAssemblyData
	HasUncheckedAssemblyData bool
	// TextInRowLimit максимальное число байт, разрешенное для текста в строке
	TextInRowLimit int
	// LargeValueTypesOutOfRow типы больших значений хранятся вне строк
	LargeValueTypesOutOfRow bool
	// IsTrackedByCDC в таблице включена система отслеживания измененных данных
	IsTrackedByCDC bool
	// LockEscalation значение параметра LOCK_ESCALATION для таблицы (TABLE | DISABLE | AUTO)
	LockEscalation string
	// IsFileTable таблица FileTable
	IsFileTable bool
	// Durability устойчивость (SCHEMA_AND_DATA | SCHEMA_ONLY)
	Durability string
	// IsMemoryOptimized memory-optimized таблица
	IsMemoryOptimized bool
	// TemporalType тип таблицы (NON_TEMPORAL_TABLE | HISTORY_TABLE | SYSTEM_VERSIONED_TEMPORAL_TABLE)
	TemporalType string
	// IsRemoteDataArchiveEnabled
	IsRemoteDataArchiveEnabled bool
	// IsExternal внешняя таблица
	IsExternal bool
	// IsNode узел графа
	IsNode bool
	// IsEdge краевая таблица графа
	IsEdge bool

	fileStreamDataSpace        sql.NullString
	historyTableSchema         sql.NullString
	historyTableName           sql.NullString
	historyRetentionPeriod     sql.NullInt32
	historyRetentionPeriodUnit sql.NullString
}

// FileStreamDataSpace возвращает наименование пространства данных для файловой группы FILESTREAM или схемы
// секционирования, состоящей из файловых групп FILESTREAM
func (table Table) FileStreamDataSpace() string {
	if table.fileStreamDataSpace.Valid {
		return table.fileStreamDataSpace.String
	}

	return ""
}

// HasHistoryRetentionPeriod проверяет, указана ли продолжительность временного периода хранения журнала
func (table Table) HasHistoryRetentionPeriod() bool {
	return table.historyRetentionPeriod.Valid
}

// HistoryRetentionPeriod возвращает продолжительность временного периода хранения журнала в единицах, указанных с
// помощью HistoryRetentionPeriodUnit()
func (table Table) HistoryRetentionPeriod() int {
	if table.historyRetentionPeriod.Valid {
		return int(table.historyRetentionPeriod.Int32)
	}

	return 0
}

// HasHistoryRetentionPeriodUnit
func (table Table) HasHistoryRetentionPeriodUnit() bool {
	return table.historyRetentionPeriodUnit.Valid
}

// HistoryRetentionPeriodUnit возвращает тип временной единицы хранения истории (INFINITE | DAY | WEEK | MONTH | YEAR)
func (table Table) HistoryRetentionPeriodUnit() string {
	if table.historyRetentionPeriodUnit.Valid {
		return table.historyRetentionPeriodUnit.String
	}

	return ""
}

// HistoryTableSchema
func (table Table) HistoryTableSchema() string {
	if table.historyTableSchema.Valid {
		return table.historyTableSchema.String
	}

	return ""
}

// HistoryTableName
func (table Table) HistoryTableName() string {
	if table.historyTableName.Valid {
		return table.historyTableName.String
	}

	return ""
}

// Tables тип коллекции таблиц в БД
type Tables map[string]*Table

// TableDefinition определение таблицы
type TableDefinition struct {
	// Table параметры таблицы
	Table *Table
	// Columns поля таблицы
	Columns Columns
	// Indexes индексы
	Indexes Indexes
	// ForeignKeys внешние ключи
	ForeignKeys ForeignKeys
	// Permissions разрешения
	Permissions UserPerms
}

// String возвращает скрипт определения таблицы. В случае возникновения ошибки при создании текста определения таблицы
// возвращает пустую строку
func (definition *TableDefinition) String() string {
	value, _ := definition.Value()
	return value
}

// Value возвращает скрипт определения таблицы. В случае возникновения ошибки при создании текста определения таблицы
// возвращает пустую строку и описание ошибки
func (definition *TableDefinition) Value() (string, error) {
	return "", nil
}

func (command *ScriptsFolderCommand) writeTableDefinition(ctx context.Context, object interface{}) (interface{}, error) {
	obj, ok := object.(IDatabaseObject)

	if !ok {
		return object, errors.New("object is not a database object")
	}

	name := obj.SchemaAndName(true)

	if obj.Type() != output.Table {
		return object, fmt.Errorf("object %s is not a table", name)
	}

	definition := &TableDefinition{
		Table:       command.tables[name],
		Columns:     command.columns[name],
		Indexes:     command.indexes[name],
		ForeignKeys: command.foreignKeys[name],
		Permissions: command.permissions[name],
	}

	value, err := definition.Value()

	if err != nil {
		return obj, err
	}

	obj.SetDefinition([]byte(value))

	return obj, nil
}

const (
	selectTables2019 = `
select tables.catalog, tables.[schema], tables.name, tables.lob_data_space, tables.lob_data_space_type,
    tables.is_default_data_space, tables.filestream_data_space, tables.lock_on_bulk_load, tables.uses_ansi_nulls,
    tables.is_replicated, tables.has_replication_filter, tables.is_merge_published, tables.is_sync_tran_subscribed,
    tables.has_unchecked_assembly_data, tables.text_in_row_limit, tables.large_value_types_out_of_row,
    tables.is_tracked_by_cdc, tables.lock_escalation, tables.is_filetable, tables.durability,
    tables.is_memory_optimized, tables.temporal_type, tables.history_table_schema, tables.history_table_name, 
	tables.is_remote_data_archive_enabled, tables.is_external, tables.history_retention_period, 
	tables.history_retention_period_unit, tables.is_node, tables.is_edge
from (
    select
        [catalog] = db_name(),
        [schema] = schema_name(objects.schema_id),
        [name] = objects.name,

        [lob_data_space] = data_spaces.name,
        [lob_data_space_type] = data_spaces.type_desc,
        [is_default_data_space] = data_spaces.is_default,
        [filestream_data_space] = filegroup_name(tables.filestream_data_space_id),
        [lock_on_bulk_load] = tables.lock_on_bulk_load,
        [uses_ansi_nulls] = tables.uses_ansi_nulls,
        [is_replicated] = tables.is_replicated,
        [has_replication_filter] = tables.has_replication_filter,
        [is_merge_published] = tables.is_merge_published,
        [is_sync_tran_subscribed] = tables.is_sync_tran_subscribed,
        [has_unchecked_assembly_data] = tables.has_unchecked_assembly_data,
        [text_in_row_limit] = tables.text_in_row_limit,
        [large_value_types_out_of_row] = tables.large_value_types_out_of_row,
        [is_tracked_by_cdc] = tables.is_tracked_by_cdc,
        [lock_escalation] = tables.lock_escalation_desc,
        [is_filetable] = tables.is_filetable,
        [durability] = tables.durability_desc,
        [is_memory_optimized] = tables.is_memory_optimized,
        [temporal_type] = tables.temporal_type_desc,
        [history_table_id] = tables.history_table_id,
        [history_table_schema] = schema_name(history_objects.schema_id),
        [history_table_name] = history_objects.name,
        [is_remote_data_archive_enabled] = tables.is_remote_data_archive_enabled,
        [is_external] = tables.is_external,
        [history_retention_period] = tables.history_retention_period,
        [history_retention_period_unit] = tables.history_retention_period_unit_desc,
        [is_node] = tables.is_node,
        [is_edge] = tables.is_edge

    from sys.tables as tables
        inner join sys.objects as objects on (tables.object_id = objects.object_id)
        left join sys.data_spaces as data_spaces on (tables.lob_data_space_id = data_spaces.data_space_id)
        left join sys.tables as history_tables
            inner join sys.objects as history_objects on (history_tables.object_id = history_objects.object_id)
        on (tables.history_table_id = history_tables.object_id)
) as tables
order by tables.catalog, tables.[schema], tables.name
`
	selectTables2017 = `
select tables.catalog, tables.[schema], tables.name, tables.lob_data_space, tables.lob_data_space_type,
    tables.is_default_data_space, tables.filestream_data_space, tables.lock_on_bulk_load, tables.uses_ansi_nulls,
    tables.is_replicated, tables.has_replication_filter, tables.is_merge_published, tables.is_sync_tran_subscribed,
    tables.has_unchecked_assembly_data, tables.text_in_row_limit, tables.large_value_types_out_of_row,
    tables.is_tracked_by_cdc, tables.lock_escalation, tables.is_filetable, tables.durability,
    tables.is_memory_optimized, tables.temporal_type, tables.history_table_schema, tables.history_table_name, 
	tables.is_remote_data_archive_enabled, tables.is_external, tables.history_retention_period, 
	tables.history_retention_period_unit, tables.is_node, tables.is_edge
from (
    select
        [catalog] = db_name(),
        [schema] = schema_name(objects.schema_id),
        [name] = objects.name,

        [lob_data_space] = data_spaces.name,
        [lob_data_space_type] = data_spaces.type_desc,
        [is_default_data_space] = data_spaces.is_default,
        [filestream_data_space] = filegroup_name(tables.filestream_data_space_id),
        [lock_on_bulk_load] = tables.lock_on_bulk_load,
        [uses_ansi_nulls] = tables.uses_ansi_nulls,
        [is_replicated] = tables.is_replicated,
        [has_replication_filter] = tables.has_replication_filter,
        [is_merge_published] = tables.is_merge_published,
        [is_sync_tran_subscribed] = tables.is_sync_tran_subscribed,
        [has_unchecked_assembly_data] = tables.has_unchecked_assembly_data,
        [text_in_row_limit] = tables.text_in_row_limit,
        [large_value_types_out_of_row] = tables.large_value_types_out_of_row,
        [is_tracked_by_cdc] = tables.is_tracked_by_cdc,
        [lock_escalation] = tables.lock_escalation_desc,
        [is_filetable] = tables.is_filetable,
        [durability] = tables.durability_desc,
        [is_memory_optimized] = tables.is_memory_optimized,
        [temporal_type] = tables.temporal_type_desc,
        [history_table_id] = tables.history_table_id,
        [history_table_schema] = schema_name(history_objects.schema_id),
        [history_table_name] = history_objects.name,
        [is_remote_data_archive_enabled] = tables.is_remote_data_archive_enabled,
        [is_external] = tables.is_external,
        [history_retention_period] = null /*tables.history_retention_period*/,
        [history_retention_period_unit] = null /*tables.history_retention_period_unit_desc*/,
        [is_node] = tables.is_node,
        [is_edge] = tables.is_edge

    from sys.tables as tables
        inner join sys.objects as objects on (tables.object_id = objects.object_id)
        left join sys.data_spaces as data_spaces on (tables.lob_data_space_id = data_spaces.data_space_id)
        left join sys.tables as history_tables
            inner join sys.objects as history_objects on (history_tables.object_id = history_objects.object_id)
        on (tables.history_table_id = history_tables.object_id)
) as tables
order by tables.catalog, tables.[schema], tables.name
`
	selectTables2016 = `
select tables.catalog, tables.[schema], tables.name, tables.lob_data_space, tables.lob_data_space_type,
    tables.is_default_data_space, tables.filestream_data_space, tables.lock_on_bulk_load, tables.uses_ansi_nulls,
    tables.is_replicated, tables.has_replication_filter, tables.is_merge_published, tables.is_sync_tran_subscribed,
    tables.has_unchecked_assembly_data, tables.text_in_row_limit, tables.large_value_types_out_of_row,
    tables.is_tracked_by_cdc, tables.lock_escalation, tables.is_filetable, tables.durability,
    tables.is_memory_optimized, tables.temporal_type, tables.history_table_schema, tables.history_table_name, 
	tables.is_remote_data_archive_enabled, tables.is_external, tables.history_retention_period, 
	tables.history_retention_period_unit, tables.is_node, tables.is_edge
from (
    select
        [catalog] = db_name(),
        [schema] = schema_name(objects.schema_id),
        [name] = objects.name,

        [lob_data_space] = data_spaces.name,
        [lob_data_space_type] = data_spaces.type_desc,
        [is_default_data_space] = data_spaces.is_default,
        [filestream_data_space] = filegroup_name(tables.filestream_data_space_id),
        [lock_on_bulk_load] = tables.lock_on_bulk_load,
        [uses_ansi_nulls] = tables.uses_ansi_nulls,
        [is_replicated] = tables.is_replicated,
        [has_replication_filter] = tables.has_replication_filter,
        [is_merge_published] = tables.is_merge_published,
        [is_sync_tran_subscribed] = tables.is_sync_tran_subscribed,
        [has_unchecked_assembly_data] = tables.has_unchecked_assembly_data,
        [text_in_row_limit] = tables.text_in_row_limit,
        [large_value_types_out_of_row] = tables.large_value_types_out_of_row,
        [is_tracked_by_cdc] = tables.is_tracked_by_cdc,
        [lock_escalation] = tables.lock_escalation_desc,
        [is_filetable] = tables.is_filetable,
        [durability] = tables.durability_desc,
        [is_memory_optimized] = tables.is_memory_optimized,
        [temporal_type] = tables.temporal_type_desc,
        [history_table_id] = tables.history_table_id,
        [history_table_schema] = schema_name(history_objects.schema_id),
        [history_table_name] = history_objects.name,
        [is_remote_data_archive_enabled] = tables.is_remote_data_archive_enabled,
        [is_external] = tables.is_external,
        [history_retention_period] = null /*tables.history_retention_period /* SQL Server 2019+ */ */,
        [history_retention_period_unit] = null /*tables.history_retention_period_unit_desc /* SQL Server 2019+ */ */,
        [is_node] = cast(0 as bit) /*tables.is_node /* SQL Server 2017+ */ */,
        [is_edge] = cast(0 as bit) /*tables.is_edge /* SQL Server 2017+ */ */

    from sys.tables as tables
        inner join sys.objects as objects on (tables.object_id = objects.object_id)
        left join sys.data_spaces as data_spaces on (tables.lob_data_space_id = data_spaces.data_space_id)
        left join sys.tables as history_tables
            inner join sys.objects as history_objects on (history_tables.object_id = history_objects.object_id)
        on (tables.history_table_id = history_tables.object_id)
) as tables
order by tables.catalog, tables.[schema], tables.name
`
)
