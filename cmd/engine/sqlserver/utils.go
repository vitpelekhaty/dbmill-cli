package sqlserver

import (
	"fmt"
	"strings"
)

// SchemaAndObject возвращает наименование объекта в формате %schema%.%name%
func SchemaAndObject(schema, objectName string) string {
	if strings.Trim(schema, " ") != "" {
		if strings.Trim(objectName, " ") != "" {
			return fmt.Sprintf("%s.%s", schema, objectName)
		}

		return schema
	} else {
		if strings.Trim(objectName, " ") != "" {
			return objectName
		}
	}

	return ""
}
