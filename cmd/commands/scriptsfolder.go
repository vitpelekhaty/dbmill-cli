package commands

import (
	"github.com/spf13/cobra"

	"github.com/vitpelekhaty/dbmill-cli/cmd/scriptsfolder"
)

// cmdScriptsFolder команда создания скриптов на основе схемы
var cmdScriptsFolder = &cobra.Command{
	Use:   "scriptsfolder",
	Short: "creates scripts based on the schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		return scriptsfolder.Run(Database, Path, IncludeData, Decrypt)
	},
}
