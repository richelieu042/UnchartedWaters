package cmd

import (
	"fmt"

	"github.com/richelieu042/chimera/v3/src/command/cobraKit"
	"github.com/spf13/cobra"
)

var versionCmd = cobraKit.NewSimpleCommand("version", "Print the version number of newApp", "", versionRun)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Println("v0.0.1")
}
