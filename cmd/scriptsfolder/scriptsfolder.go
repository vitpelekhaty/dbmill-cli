package scriptsfolder

import (
	"github.com/vitpelekhaty/dbmill-cli/cmd/engine"
)

// Run создает скрипты на основе схемы
func Run(connection string, path string, includeData, decrypt bool, options ...engine.DatabaseOption) error {
	n, err := engine.NewDatabaseConnection(connection, options...)

	if err != nil {
		return err
	}

	return n.ScriptsFolder(path, includeData, decrypt)
}
