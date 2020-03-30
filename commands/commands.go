package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cmdRoot команда - точка входа в приложение
var cmdRoot = &cobra.Command{
	Use:   "dbmill-cli",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

// Execute выполняет команду приложения
func Execute() error {
	return cmdRoot.Execute()
}

func init() {
	cmdScriptsFolder.Flags().StringVarP(&Path, "path", "d", "",
		"path to the directory where scripts will be created")
	cmdScriptsFolder.Flags().StringVarP(&Database, "db", "D", "",
		"database to create scripts")
	cmdScriptsFolder.Flags().StringVarP(&DirStructFilename, "dir-struct", "S", "",
		"path to a file that describes a directory structure where the scripts will be created")
	cmdScriptsFolder.Flags().StringVarP(&LogFilename, "log", "l", "",
		"path to a log file")
	cmdScriptsFolder.Flags().StringVarP(&LogLevel, "log-level", "L", "info",
		"log level: trace, debug, info (default), warning, error, fatal, panic")
	cmdScriptsFolder.Flags().StringVarP(&FilterPath, "filter-path", "F", "",
		"path to a file that contains a list of objects for which scripts will be created; replaces --filter "+
			"if it is empty")
	cmdScriptsFolder.Flags().StringVarP(&Username, "username", "U", "",
		"database username; replaces a username listed in a database connection string")
	cmdScriptsFolder.Flags().StringVarP(&Password, "password", "P", "",
		"database user password; replaces a password listed in a database connection string")

	cmdScriptsFolder.Flags().StringArrayVarP(&Filter, "filter", "f", nil,
		"names of objects for which scripts will be created; regular expressions are permissible; "+
			"scripts will be created for all objects if the option is empty; replaces --filter-path")

	cmdScriptsFolder.Flags().BoolVarP(&Decrypt, "decrypt", "", false,
		"decrypt objects")
	cmdScriptsFolder.Flags().BoolVarP(&IncludeData, "include-data", "", false,
		"save data in scripts")

	cmdRoot.AddCommand(cmdScriptsFolder, cmdVersion)
}
