package cmd

import (
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/command/cobraKit"
	"github.com/richelieu042/chimera/v3/src/log/console"
	"github.com/richelieu042/chimera/v3/src/log/zapKit"
	"github.com/spf13/cobra"
)

var cleanCmd = cobraKit.NewSimpleCommand("clean", "清理 adb 环境", "", cleanRun)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func cleanRun(cmd *cobra.Command, args []string) {
	success := checkAdbEnv()
	if !success {
		return
	}

	logger := zapKit.NewSimpleConsoleLogger()
	if err := adbKit.Clean(logger); err != nil {
		console.Errorf("Clean failed, error %s", err)
		return
	}
	console.Infoln("Clean successfully.")
}
