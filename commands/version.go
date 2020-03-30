package commands

import (
	"github.com/spf13/cobra"
)

var (
	// GitBranch ветка репозитория исходного кода, из которого собрано приложение
	GitBranch string
	// GitCommit идентификатор фиксации исходного кода, из которого собрано приложение
	GitCommit string
	// Время сборки приложения
	Built string
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf()
	},
}
