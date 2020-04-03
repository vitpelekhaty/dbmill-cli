package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine"
	"github.com/vitpelekhaty/dbmill-cli/cmd/input"
	"github.com/vitpelekhaty/dbmill-cli/cmd/scriptsfolder"
)

// cmdScriptsFolder команда создания скриптов на основе схемы
var cmdScriptsFolder = &cobra.Command{
	Use:   "scriptsfolder",
	Short: "creates scripts based on the schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.Trim(Username, " ") == "" || strings.Trim(Password, " ") == "" {
			user, pwd, err := engine.Credentials(Database)

			if err != nil {
				return err
			}

			if strings.Trim(Username, " ") == "" {
				Username = user
			}

			if strings.Trim(Password, " ") == "" {
				Password = pwd
			}
		}

		if strings.Trim(Username, " ") == "" {
			Username = input.Username()
		}

		if strings.Trim(Password, " ") == "" {
			Password = input.Password()
		}

		return scriptsfolder.Run(Database, Path, IncludeData, Decrypt)
	},
}
