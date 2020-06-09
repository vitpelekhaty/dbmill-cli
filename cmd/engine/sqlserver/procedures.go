package sqlserver

import (
	"context"
	"errors"
	"fmt"
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
		return object, fmt.Errorf("object %s is not a procedure", obj.SchemaAndName())
	}

	definition := string(obj.Definition())

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

	obj.SetDefinition([]byte(definition))

	return obj, nil
}
