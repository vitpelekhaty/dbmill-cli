package sqlserver

import (
	"context"
	"fmt"
	"sort"
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

	objectName := obj.Name()
	permissions := command.permissions[objectName]

	var stringPermissions string

	for user, userPermissions := range permissions {
		for state, permissions := range userPermissions {
			if len(permissions) > 0 {
				perms := permissions.List()
				sort.Strings(perms)

				stringPermissions = strings.Join(perms, ",\n  ")
				stringPermissions = fmt.Sprintf("%s\n  %s\nON SCHEMA :: [%s] TO [%s]\nGO", state.String(),
					stringPermissions, objectName, user)

				definition = fmt.Sprintf("%s\n\n%s", definition, stringPermissions)
			}
		}
	}

	obj.SetDefinition([]byte(definition))

	return obj, nil
}

const schemaDefinition = `CREATE SCHEMA [%s] AUTHORIZATION [%s]
GO`

const schemaShortDefinition = `CREATE SCHEMA [%s]
GO`
