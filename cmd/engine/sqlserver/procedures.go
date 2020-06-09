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

	definition := string(obj.Definition())

	if strings.Trim(definition, " ") == "" {
		return object, nil
	}

	definition = fmt.Sprintf("%s\nGO", strings.Trim(definition, "\n"))

	mod := make([]string, 0)

	if obj.QuotedIdentifierValid() {
		if obj.QuotedIdentifier() {
			mod = append(mod, "QUOTED_IDENTIFIER")
		}
	}

	if obj.ANSINullsValid() {
		if obj.ANSINulls() {
			mod = append(mod, "ANSI_NULLS")
		}
	}

	s := strings.Join(mod, ", ")

	if strings.Trim(s, " ") != "" && strings.Trim(definition, " ") != "" {
		definition = fmt.Sprintf("SET %s ON\nGO\n\n%s", s, definition)
	}

	name := obj.SchemaAndName(true)
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

	obj.SetDefinition([]byte(definition))

	return obj, nil
}
