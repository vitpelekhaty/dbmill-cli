package commands

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	// Version ветка репозитория исходного кода, из которого собрано приложение
	Version string
	// GitCommit идентификатор фиксации исходного кода, из которого собрано приложение
	GitCommit string
	// Время сборки приложения
	Built string
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n%s\n%s\n\n", mill, app, appDescription)

		tw := tabwriter.NewWriter(os.Stdout, 0, 20, 0, '\t', 0)
		fmt.Fprintf(tw, "GOARCH\t%s\t\nGOOS\t%s\t\nVersion\t%s\t\nCommit\t%s\t\nBuilt\t%s\t\n",
			runtime.GOARCH, runtime.GOOS, Version, GitCommit, Built)
		tw.Flush()
	},
}

const mill = `
      ##        ##
       ##      ##
        ##|--|##
       | ##  ## |
      |___####___|
     |~~~~####~~~~|
      ===##==##===
     |..##....##..|
    |..##......##..|
   |..##........##..|
  |........__........|
 |........|__|........|
|_________|__|_________|
`

const app = "dbmill-cli"
const appDescription = ""
