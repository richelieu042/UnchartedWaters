package cmd

import (
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/command/cobraKit"
	"github.com/richelieu042/chimera/v3/src/log/console"
	"github.com/spf13/cobra"
)

var checkCmd = cobraKit.NewSimpleCommand("check", "检查 adb 环境", "", checkRun)

func init() {
	rootCmd.AddCommand(checkCmd)
}

func checkRun(cmd *cobra.Command, args []string) {
	checkAdbEnv()
}

func checkAdbEnv() (success bool) {
	path, version, err := adbKit.Check()
	if err != nil {
		console.Errorf("Check failed, error: %s", err)
		return
	}
	console.Infof("adb path: %s", path)
	console.Infof("adb version: \n%s", version)
	console.Info("Check successfully.")
	success = true
	return
}
