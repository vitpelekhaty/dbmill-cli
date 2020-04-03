package engine

import (
	"errors"
)

// ErrorUnsupportedDatabaseType ошибка "Неподдерживаемая СУБД"
var ErrorUnsupportedDatabaseType = errors.New("unsupported database type")
