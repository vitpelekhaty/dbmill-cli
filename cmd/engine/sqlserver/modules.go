package sqlserver

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

func (command *ScriptsFolderCommand) writeProcedureDefinition(ctx context.Context, object interface{}) (interface{},
	error) {
	obj, ok := object.(ISQLModule)

	if !ok {
		return object, errors.New("object is not a SQL module")
	}

	if obj.Type() != output.Procedure {
		return object, fmt.Errorf("object %s is not a procedure", obj.SchemaAndName(true))
	}

	return command.writeModuleDefinition(ctx, obj)
}

func (command *ScriptsFolderCommand) writeFunctionDefinition(ctx context.Context, object interface{}) (interface{},
	error) {
	obj, ok := object.(ISQLModule)

	if !ok {
		return object, errors.New("object is not a SQL module")
	}

	if obj.Type() != output.Function {
		return object, fmt.Errorf("object %s is not a function", obj.SchemaAndName(true))
	}

	return command.writeModuleDefinition(ctx, obj)
}

func (command *ScriptsFolderCommand) writeViewDefinition(ctx context.Context, object interface{}) (interface{},
	error) {
	obj, ok := object.(ISQLModule)

	if !ok {
		return object, errors.New("object is not a SQL module")
	}

	if obj.Type() != output.View {
		return object, fmt.Errorf("object %s is not a view", obj.SchemaAndName(true))
	}

	return command.writeModuleDefinition(ctx, obj)
}

func (command *ScriptsFolderCommand) writeTriggerDefinition(ctx context.Context, object interface{}) (interface{},
	error) {
	obj, ok := object.(ISQLModule)

	if !ok {
		return object, errors.New("object is not a SQL module")
	}

	if obj.Type() != output.Trigger {
		return object, fmt.Errorf("object %s is not a DDL trigger", obj.SchemaAndName(true))
	}

	return command.writeModuleDefinition(ctx, obj)
}

func (command *ScriptsFolderCommand) writeModuleDefinition(ctx context.Context, object ISQLModule) (ISQLModule, error) {
	definition := string(object.Definition())

	if strings.Trim(definition, " ") == "" {
		return object, nil
	}

	definition = fmt.Sprintf("%s\nGO", strings.Trim(definition, "\n"))

	mod := make([]string, 0)

	if object.QuotedIdentifierValid() {
		if object.QuotedIdentifier() {
			mod = append(mod, "QUOTED_IDENTIFIER")
		}
	}

	if object.ANSINullsValid() {
		if object.ANSINulls() {
			mod = append(mod, "ANSI_NULLS")
		}
	}

	s := strings.Join(mod, ", ")

	if strings.Trim(s, " ") != "" && strings.Trim(definition, " ") != "" {
		definition = fmt.Sprintf("SET %s ON\nGO\n%s", s, definition)
	}

	if !command.skipPermissions {
		name := object.SchemaAndName(true)
		permissions := command.permissions[name]

		if len(permissions) > 0 {
			users := permissions.Users()
			sort.Strings(users)

			for _, user := range users {
				states := permissions[user]

				for state, perms := range states {
					ps := perms.List()
					sort.Strings(ps)

					for _, p := range ps {
						definition = fmt.Sprintf("%s\n\n%s %s ON %s TO [%s]\nGO", definition, state, p, name, user)
					}
				}
			}
		}
	}

	object.SetDefinition([]byte(definition))

	return object, nil
}
