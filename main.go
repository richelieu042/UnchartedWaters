package main

import (
	"os"

	"github.com/richelieu-yang/UnchartedWaters/src/cmd"
	"github.com/richelieu042/chimera/v3/src/core/sliceKit"
	"github.com/richelieu042/chimera/v3/src/log/console"
)

func main() {
	/* test */
	{
		//addr := "192.168.60.205:16384"
		addr := "127.0.0.1:5555"
		//addr := "127.0.0.1:5585"

		os.Args = []string{"uw", "--addr", addr, "--clean"}

		verbose := true
		disableSail := false
		disableBattle := true
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
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
