package cmd

import (
	"fmt"

	"github.com/richelieu042/chimera/v3/src/command/cobraKit"
	"github.com/spf13/cobra"
)

var (
	rootCmd = cobraKit.NewSimpleCommand("uw", "《大航海时代：传说》的adb脚本。", "", rootRun)
)

func rootRun(cmd *cobra.Command, args []string) {
	fmt.Println("$$$")
}

func Execute() error {
	return rootCmd.Execute()
}
