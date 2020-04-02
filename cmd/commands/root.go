package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cmdRoot команда - точка входа в приложение
var cmdRoot = &cobra.Command{
	Use:   app,
	Short: appDescription,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}
