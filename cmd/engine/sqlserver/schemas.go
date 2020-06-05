package sqlserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

func (command *ScriptsFolderCommand) writeSchemaDefinition(ctx context.Context, object interface{}) (interface{},
	error) {
	obj := object.(databaseObject)

	if obj.Type() != output.Schema {
		return object, fmt.Errorf("object %s is not a schema", obj.SchemaAndName())
	}

	owner := obj.Owner()

	var definition string

	if strings.Trim(owner, " ") != "" {
		definition = fmt.Sprintf(schemaDefinition, obj.Schema(), owner)
	} else {
		definition = fmt.Sprintf(schemaShortDefinition, obj.Schema())
	}

	obj.SetDefinition([]byte(definition))

	return obj, nil
}

const schemaDefinition = `CREATE SCHEMA [%s] AUTHORIZATION [%s]
GO`

const schemaShortDefinition = `CREATE SCHEMA [%s]
GO`
