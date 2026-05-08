package main

import (
	"os"
	"time"

	"github.com/richelieu-yang/UnchartedWaters/src/cmd"
	"github.com/richelieu-yang/UnchartedWaters/src/conf"
	"github.com/richelieu042/chimera/v3/src/core/sliceKit"
	"github.com/richelieu042/chimera/v3/src/log/console"
	"go.uber.org/zap"
)

var (
	addr                                string
	verbose, disableSail, disableBattle bool
)

func main() {
	mode := 0

	switch mode {
	case 0: /* 本地 - 跑商 */
		addr = "127.0.0.1:5555"
		//addr = "127.0.0.1:5585"

		verbose = true
		disableSail = false
		disableBattle = true
	case 1: /* 远程 - 跑商 */
		addr = "192.168.60.205:16384"

		verbose = true
		disableSail = false
		disableBattle = true
		conf.SetDefSleepInterval(time.Millisecond * 10)
	case 2: /* 本地 - 海战 */
		//addr = "192.168.60.205:16384"
		addr = "127.0.0.1:5555"
		//addr = "127.0.0.1:5585"

		verbose = true
		disableSail = true
		disableBattle = false
	default:
		console.Fatal("Unknown mode.", zap.Int("mode", mode))
	}

	os.Args = []string{"uw", "--addr", addr, "--clean"}
	if verbose {
		os.Args = append(os.Args, "--verbose")
	}
	if disableSail {
		os.Args = append(os.Args, "--disable_sail")
	}
	if disableBattle {
		os.Args = append(os.Args, "--disable_battle")
	}
	console.Warnf("command: %s", sliceKit.Join(os.Args, " "))

	//os.Args = []string{"uw", "clean"}
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
