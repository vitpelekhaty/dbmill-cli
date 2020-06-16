package sqlserver

import (
	"fmt"
	"strings"
)

// SchemaAndObject возвращает наименование объекта в формате %Schema%.%name%
func SchemaAndObject(schema, objectName string, useBrackets bool) string {
	if strings.Trim(schema, " ") != "" {
		if strings.Trim(objectName, " ") != "" {
			if useBrackets {
				return fmt.Sprintf("[%s].[%s]", schema, objectName)
			} else {
				return fmt.Sprintf("%s.%s", schema, objectName)
			}
		}

		if useBrackets {
			return "[" + schema + "]"
		} else {
			return schema
		}
	} else {
		if strings.Trim(objectName, " ") != "" {
			if useBrackets {
				return "[" + objectName + "]"
			} else {
				return objectName
			}
		}
	}

	return ""
}
