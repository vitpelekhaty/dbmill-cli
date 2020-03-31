package commands

import (
	"github.com/spf13/cobra"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine"
)

// cmdScriptsFolder команда создания скриптов на основе схемы
var cmdScriptsFolder = &cobra.Command{
	Use:   "scriptsfolder",
	Short: "creates scripts based on the schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ScriptsFolder()
	},
}

// ScriptsFolder создает скрипты на основе схемы
func ScriptsFolder() error {
	n, err := engine.NewDatabaseConnection(Database)

	if err != nil {
		return err
	}

	return n.ScriptsFolder(Path, IncludeData, Decrypt)
}
