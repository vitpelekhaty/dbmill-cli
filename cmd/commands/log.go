package commands

import (
	"fmt"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

// ParseLogLevel возвращает уровень лога по его наименованию
func ParseLogLevel(level string) (log.Level, error) {
	switch level {
	case "trace":
		return log.TraceLevel, nil
	case "debug":
		return log.DebugLevel, nil
	case "info":
		return log.InfoLevel, nil
	case "warning":
		return log.WarningLevel, nil
	case "error":
		return log.ErrorLevel, nil
	case "fatal":
		return log.FatalLevel, nil
	case "panic":
		return log.PanicLevel, nil
	default:
		return log.InfoLevel, fmt.Errorf("unknown log level %s", level)
	}
}
