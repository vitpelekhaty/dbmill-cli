package sqlserver

import (
	"context"
	"database/sql"
)

// MetadataReader объект чтения метаданных из базы
type MetadataReader struct {
	db *sql.DB

	serverVersion int
}

// NewMetadataReader конструктор MetadataReader
func NewMetadataReader(engine *Engine, serverVersion int) (*MetadataReader, error) {
	return &MetadataReader{db: engine.db, serverVersion: serverVersion}, nil
}

// Permissions возвращает разрешения на объекты БД
func (reader *MetadataReader) Permissions(ctx context.Context) (ObjectPermissions, error) {
	stmt, err := reader.db.PrepareContext(ctx, selectPermissions)

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
		nullableSchema sql.NullString
		schema         string
		object         string
		permission     string
		state          string
		user           string
	)

	perms := make(ObjectPermissions)

	for rows.Next() {
		err = rows.Scan(&nullableSchema, &object, &permission, &state,
			&user)

		if err != nil {
			return nil, err
		}

		if nullableSchema.Valid {
			schema = nullableSchema.String
		}

		err = perms.Append(schema, object, permission, state, user)

		if err != nil {
			return nil, err
		}
	}

	return perms, nil
}

// UserDefinedTypes возвращает список пользовательских типов
func (reader *MetadataReader) UserDefinedTypes(ctx context.Context) (UserDefinedTypes, error) {
	stmt, err := reader.db.PrepareContext(ctx, selectUserDefinedTypes)

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
		parentTypeName    sql.NullString
		maxLength         sql.NullString
		precision         sql.NullInt32
		scale             sql.NullInt32
		collation         sql.NullString
		isNullable        bool
		isTableType       bool
		isMemoryOptimized bool
	)

	types := make(UserDefinedTypes)

	for rows.Next() {
		err = rows.Scan(&catalog, &typeName, &schema, &parentTypeName, &maxLength, &precision, &scale, &collation,
			&isNullable, &isTableType, &isMemoryOptimized)

		if err != nil {
			return nil, err
		}

		t := &UserDefinedType{
			Catalog:           catalog,
			TypeName:          typeName,
			Schema:            schema,
			parentTypeName:    parentTypeName,
			maxLength:         maxLength,
			precision:         precision,
			scale:             scale,
			collation:         collation,
			IsNullable:        isNullable,
			IsTableType:       isTableType,
			IsMemoryOptimized: isMemoryOptimized,
		}

		types.append(t)
	}

	return types, nil
}

// ObjectColumns возвращает поля объектов БД (таблиц/табличных типов)
func (reader *MetadataReader) ObjectColumns(ctx context.Context) (ObjectColumns, error) {
	stmt, err := reader.db.PrepareContext(ctx, selectColumns)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns := make(ObjectColumns)

	var (
		catalog                       string
		schema                        string
		objectName                    string
		columnID                      int
		columnName                    string
		description                   sql.NullString
		typeName                      string
		typeSchema                    string
		isUserDefinedType             bool
		maxLength                     sql.NullString
		precision                     sql.NullInt32
		scale                         sql.NullInt32
		collation                     sql.NullString
		isNullable                    bool
		isANSIPadded                  bool
		isRowGUIDCol                  bool
		isIdentity                    bool
		seedValue                     sql.NullInt32
		incValue                      sql.NullInt32
		isComputed                    bool
		isPersisted                   sql.NullBool
		compute                       sql.NullString
		isFileStream                  bool
		isReplicated                  bool
		isNonSQLSubscribed            bool
		isMergePublished              bool
		isDTSReplicated               bool
		isXMLDocument                 bool
		xmlSchemaCollectionSchemaName sql.NullString
		xmlSchemaCollectionName       sql.NullString
		defaultConstraint             sql.NullString
		def                           sql.NullString
		isSparse                      bool
		isColumnSet                   bool
		generateAlways                sql.NullString
		isHidden                      bool
		isMasked                      bool
		maskingFunc                   sql.NullString
		encryptionKey                 sql.NullString
		encryptionKeyDatabaseName     sql.NullString
		encryptionAlgorithm           sql.NullString
		encryptionType                sql.NullString
	)

	for rows.Next() {
		err = rows.Scan(&catalog, &schema, &objectName, &columnID, &columnName, &description, &typeName, &typeSchema,
			&isUserDefinedType, &maxLength, &precision, &scale, &collation, &isNullable, &isANSIPadded, &isRowGUIDCol,
			&isIdentity, &seedValue, &incValue, &isComputed, &isPersisted, &compute, &isFileStream, &isReplicated,
			&isNonSQLSubscribed, &isMergePublished, &isDTSReplicated, &isXMLDocument, &xmlSchemaCollectionSchemaName,
			&xmlSchemaCollectionName, &defaultConstraint, &def, &isSparse, &isColumnSet, &generateAlways, &isHidden,
			&isMasked, &maskingFunc, &encryptionKey, &encryptionType, &encryptionAlgorithm, &encryptionKeyDatabaseName)

		if err != nil {
			return nil, err
		}

		column := &Column{
			ID:                            columnID,
			Name:                          columnName,
			description:                   description,
			TypeName:                      typeName,
			TypeSchema:                    typeSchema,
			IsUserDefinedType:             isUserDefinedType,
			maxLength:                     maxLength,
			precision:                     precision,
			scale:                         scale,
			collation:                     collation,
			IsNullable:                    isNullable,
			IsANSIPadded:                  isANSIPadded,
			IsRowGUIDCol:                  isRowGUIDCol,
			IsIdentity:                    isIdentity,
			identitySeedValue:             seedValue,
			identityIncrementValue:        incValue,
			isComputed:                    isComputed,
			isPersisted:                   isPersisted,
			compute:                       compute,
			IsFileStream:                  isFileStream,
			IsReplicated:                  isReplicated,
			IsNonSQLSubscribed:            isNonSQLSubscribed,
			IsMergePublished:              isMergePublished,
			IsDTSReplicated:               isDTSReplicated,
			IsXMLDocument:                 isXMLDocument,
			xmlSchemaCollectionSchemaName: xmlSchemaCollectionSchemaName,
			xmlSchemaCollectionName:       xmlSchemaCollectionName,
			defaultConstraint:             defaultConstraint,
			defaultConstraintDefinition:   def,
			IsSparse:                      isSparse,
			IsColumnSet:                   isColumnSet,
			generateAlways:                generateAlways,
			IsHidden:                      isHidden,
			IsMasked:                      isMasked,
			maskingFunction:               maskingFunc,
			encryptionKey:                 encryptionKey,
			encryptionType:                encryptionType,
			encryptionAlgorithm:           encryptionAlgorithm,
			encryptionKeyDatabaseName:     encryptionKeyDatabaseName,
		}

		err = columns.Append(schema, objectName, column)

		if err != nil {
			return nil, err
		}
	}

	return columns, nil
}

// DatabaseCollation возвращает collation базы данных
func (meta *MetadataReader) DatabaseCollation(ctx context.Context) (string, error) {
	stmt, err := meta.db.PrepareContext(ctx, selectDatabaseCollation)

	if err != nil {
		return "", nil
	}

	defer stmt.Close()

	var collation string

	err = stmt.QueryRowContext(ctx).Scan(&collation)

	return collation, err
}

var selectIndexesQueries = map[int]string{
	13: selectIndexes2016,
	15: selectIndexes2019,
}

// selectIndexesQuery возвращает текст запроса набора индексов для соответствующей версии SQL Server.
// Если для указанной версии нет варианта текста запроса, то возвращается текст для минимальной поддерживаемой
// версии
func (meta *MetadataReader) selectIndexesQuery() string {
	if query, ok := selectIndexesQueries[meta.serverVersion]; ok {
		return query
	}

	return selectIndexes2016
}

// ObjectsIndexes возвращает справочник индексов из БД, сгруппированных по объектам
func (meta *MetadataReader) Indexes(ctx context.Context) (ObjectsIndexes, error) {
	stmt, err := meta.db.PrepareContext(ctx, meta.selectIndexesQuery())

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	indexes := make(ObjectsIndexes)

	var (
		catalog                  string
		schema                   string
		objectName               string
		indexName                string
		indexType                string
		isUnique                 bool
		isPrimaryKey             bool
		isUniqueConstraint       bool
		ignoreDupKey             bool
		fillFactor               byte
		isPadded                 bool
		isDisabled               bool
		isHypothetical           bool
		isIgnoredInOptimization  bool
		allowRowLocks            bool
		allowPageLocks           bool
		suppressDupKeyMessages   bool
		autoCreated              bool
		optimizeForSequentialKey bool
		hasFilter                bool
		filterDefinition         sql.NullString
		indexColumnID            int
		columnName               string
		isDescendingKey          bool
		isIncludedColumn         bool
		keyOrdinal               int
		partitionOrdinal         int
		columnStoreOrderOrdinal  int
		bucketCount              sql.NullInt64
		description              sql.NullString

		name string
	)

	for rows.Next() {
		err = rows.Scan(&catalog, &schema, &objectName, &indexName, &indexType, &isUnique, &isPrimaryKey,
			&isUniqueConstraint, &ignoreDupKey, &fillFactor, &isPadded, &isDisabled, &isHypothetical,
			&isIgnoredInOptimization, &allowRowLocks, &allowPageLocks, &suppressDupKeyMessages, &autoCreated,
			&optimizeForSequentialKey, &hasFilter, &filterDefinition, &indexColumnID, &columnName, &isDescendingKey,
			&isIncludedColumn, &keyOrdinal, &partitionOrdinal, &columnStoreOrderOrdinal, &bucketCount, &description)

		if err != nil {
			return nil, err
		}

		name = SchemaAndObject(schema, objectName, true)

		column := &IndexedColumn{
			ID:                      indexColumnID,
			Name:                    columnName,
			IsDescendingKey:         isDescendingKey,
			KeyOrdinal:              keyOrdinal,
			PartitionOrdinal:        partitionOrdinal,
			ColumnStoreOrderOrdinal: columnStoreOrderOrdinal,
		}

		if indexes[name] == nil {
			indexes[name] = make(Indexes)
		}

		if index, ok := indexes[name][indexName]; ok {
			if isIncludedColumn {
				if _, ok := index.IncludedColumns[columnName]; !ok {
					index.IncludedColumns[columnName] = column
				}
			} else {
				if _, ok := index.Columns[columnName]; !ok {
					index.Columns[columnName] = column
				}
			}
		} else {
			index := &Index{
				Name:                     indexName,
				Type:                     indexType,
				IsUnique:                 isUnique,
				IsPrimaryKey:             isPrimaryKey,
				IsUniqueConstraint:       isUniqueConstraint,
				IgnoreDupKey:             ignoreDupKey,
				FillFactor:               fillFactor,
				IsPadded:                 isPadded,
				IsDisabled:               isDisabled,
				IsHypothetical:           isHypothetical,
				IsIgnoredInOptimization:  isIgnoredInOptimization,
				AllowRowLocks:            allowRowLocks,
				AllowPageLocks:           allowPageLocks,
				SuppressDupKeyMessages:   suppressDupKeyMessages,
				AutoCreated:              autoCreated,
				OptimizeForSequentialKey: optimizeForSequentialKey,
				Columns:                  make(IndexedColumns),
				IncludedColumns:          make(IndexedColumns),
				hasFilter:                hasFilter,
				filterDefinition:         filterDefinition,
				bucketCount:              bucketCount,
				description:              description,
			}

			if isIncludedColumn {
				index.IncludedColumns[columnName] = column
			} else {
				index.Columns[columnName] = column
			}

			indexes[name][indexName] = index
		}
	}

	return indexes, nil
}

// ObjectsForeignKeys возвращает справочник внешних ключей
func (meta *MetadataReader) ForeignKeys(ctx context.Context) (ObjectsForeignKeys, error) {
	stmt, err := meta.db.PrepareContext(ctx, selectForeignKeys)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	foreignKeys := make(ObjectsForeignKeys)

	var (
		catalog                 string
		foreignKeyName          string
		constraintColumnID      int
		parentObjectSchema      string
		parentObjectName        string
		parentColumnName        string
		referencedObjectSchema  string
		referencedObjectName    string
		referencedColumnsName   string
		isDisabled              bool
		isNotForReplication     bool
		isNotTrusted            bool
		deleteReferentialAction string
		updateReferentialAction string
		description             sql.NullString

		name string
	)

	for rows.Next() {
		err = rows.Scan(&catalog, &foreignKeyName, &constraintColumnID, &parentObjectSchema, &parentObjectName,
			&parentColumnName, &referencedObjectSchema, &referencedObjectName, &referencedColumnsName, &isDisabled,
			&isNotForReplication, &isNotTrusted, &deleteReferentialAction, &updateReferentialAction, &description)

		if err != nil {
			return nil, err
		}

		name = SchemaAndObject(parentObjectSchema, parentObjectName, true)

		columnReference := &ColumnReference{
			ID:               constraintColumnID,
			Column:           parentColumnName,
			ReferencedColumn: referencedColumnsName,
		}

		if foreignKeys[name] == nil {
			foreignKeys[name] = make(ForeignKeys)
		}

		if fk, ok := foreignKeys[name][foreignKeyName]; ok {
			if _, ok := fk.ColumnsReferences[parentColumnName]; !ok {
				fk.ColumnsReferences[parentColumnName] = columnReference
			}
		} else {
			fk := &ForeignKey{
				Name:                    foreignKeyName,
				ReferencedObjectSchema:  referencedObjectSchema,
				ReferencedObjectName:    referencedObjectName,
				IsDisabled:              isDisabled,
				IsNotForReplication:     isNotForReplication,
				IsNotTrusted:            isNotTrusted,
				DeleteReferentialAction: deleteReferentialAction,
				UpdateReferentialAction: updateReferentialAction,
				ColumnsReferences:       make(map[string]*ColumnReference),
				description:             description,
			}

			fk.ColumnsReferences[parentColumnName] = columnReference

			foreignKeys[name][foreignKeyName] = fk
		}
	}

	return foreignKeys, nil
}

var selectTablesQueries = map[int]string{
	13: selectTables2016,
	14: selectTables2017,
	15: selectTables2019,
}

// selectTablesQuery возвращает текст запроса набора таблиц для соответствующей версии SQL Server.
// Если для указанной версии нет варианта текста запроса, то возвращается текст для минимальной поддерживаемой
// версии
func (meta *MetadataReader) selectTablesQuery() string {
	if query, ok := selectTablesQueries[meta.serverVersion]; ok {
		return query
	}

	return selectTables2016
}

// Tables возвращает коллекцию пользовтельских таблиц, имеющихся в БД
func (meta *MetadataReader) Tables(ctx context.Context) (Tables, error) {
	stmt, err := meta.db.PrepareContext(ctx, meta.selectTablesQuery())

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	return nil, nil
}
