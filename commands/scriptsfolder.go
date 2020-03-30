package commands

import (
	"github.com/spf13/cobra"
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
	return nil
}
