package sqlserver

import (
	"context"
	"database/sql"
)

// MetadataReader объект чтения метаданных из базы
type MetadataReader struct {
	db *sql.DB
}

// NewMetadataReader конструктор MetadataReader
func NewMetadataReader(engine *Engine) (*MetadataReader, error) {
	return &MetadataReader{db: engine.db}, nil
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
		parentTypeName    string
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
			ParentTypeName:    parentTypeName,
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
