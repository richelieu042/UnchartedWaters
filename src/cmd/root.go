package cmd

import (
	"github.com/richelieu-yang/UnchartedWaters/src/log"
	"github.com/richelieu-yang/UnchartedWaters/src/logic"
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/command/cobraKit"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	addr    string
	clean   bool
	verbose bool

	rootCmd = cobraKit.NewSimpleCommand("uw", "《大航海时代：传说》的 adb 脚本。", "", rootRun)
)

func init() {
	rootCmd.Flags().StringVarP(&addr, "addr", "", "127.0.0.1:5555", "adb连接地址")
	//if err := rootCmd.MarkFlagRequired("addr"); err != nil {
	//	panic(err)
	//}

	rootCmd.Flags().BoolVarP(&clean, "clean", "", false, "在adb连接前，清理adb环境")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "", false, "更多的输出")
}

func rootRun(cmd *cobra.Command, args []string) {
	var level zapcore.Level
	if verbose {
		level = zap.DebugLevel
	} else {
		level = zap.InfoLevel
	}
	log.SetLevel(level)

	logger := log.NewLogger("")
	logger.Sugar().Infof("addr: [%s]", addr)
	logger.Sugar().Infof("clean: [%t]", clean)
	logger.Sugar().Infof("verbose: [%t]", verbose)

	client, err := adbKit.NewClient(addr, clean, logger)
	if err != nil {
		logger.Sugar().Panicf("Fail to new adb client, error: %+v", err)
		return
	}
	logic.Start(client, logger)
}

func Execute() error {
	return rootCmd.Execute()
}
