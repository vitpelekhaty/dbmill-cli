package sqlserver

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

func (command *ScriptsFolderCommand) writeSchemaDefinition(ctx context.Context, object interface{}) (interface{},
	error) {
	obj, ok := object.(IDatabaseObject)

	if !ok {
		return object, errors.New("object is not a database object")
	}

	if obj.Type() != output.Schema {
		return object, fmt.Errorf("object %s is not a Schema", obj.SchemaAndName(true))
	}

	owner := obj.Owner()

	var definition string

	if strings.Trim(owner, " ") != "" {
		definition = fmt.Sprintf(schemaDefinition, obj.Schema(), owner)
	} else {
		definition = fmt.Sprintf(schemaShortDefinition, obj.Schema())
	}

	if !command.skipPermissions {
		objectName := obj.SchemaAndName(true)
		permissions := command.permissions[objectName]

		if len(permissions) > 0 {
			var stringPermissions string

			for user, userPermissions := range permissions {
				for state, permissions := range userPermissions {
					if len(permissions) > 0 {
						perms := permissions.List()
						sort.Strings(perms)

						stringPermissions = strings.Join(perms, ",\n  ")
						stringPermissions = fmt.Sprintf("%s\n  %s\nON SCHEMA :: %s TO [%s]\nGO", state.String(),
							stringPermissions, objectName, user)

						definition = fmt.Sprintf("%s\n\n%s", definition, stringPermissions)
					}
				}
			}
		}
	}

	description := obj.Description()

	if strings.Trim(description, " ") != "" {
		description = fmt.Sprintf(schemaDescription, obj.SchemaAndName(false), description)
		definition = fmt.Sprintf("%s\n%s", definition, description)
	}

	obj.SetDefinition([]byte(definition))

	return obj, nil
}

const schemaDefinition = `CREATE SCHEMA [%s] AUTHORIZATION [%s]
GO`

const schemaShortDefinition = `CREATE SCHEMA [%s]
GO`

const schemaDescription = `
EXECUTE sp_addextendedproperty @name = N'MS_Description', @level0type = N'SCHEMA', @level0name = N'%s', @value = N'%s'
GO`
