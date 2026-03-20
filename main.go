package main

import (
	"os"

	"github.com/richelieu-yang/UnchartedWaters/src/cmd"
)

func main() {
	/* test */
	{
		//addr := "192.168.60.205:5555"
		addr := "192.168.60.205:16384"
		//addr := "127.0.0.1:5555"
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
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
