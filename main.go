package main

import (
	"os"

	"github.com/richelieu-yang/UnchartedWaters/src/cmd"
)

func main() {
	/* test */
	{
		addr := "127.0.0.1:5555"
		//addr := "127.0.0.1:5585"

		//os.Args = []string{"uw", "--addr", addr, "--clean", "--verbose", "--disable_battle"}
		os.Args = []string{"uw", "--addr", addr, "--clean", "--verbose", "--disable_sail"}
		//os.Args = []string{"uw", "--addr", addr, "--clean", "--verbose"}
		//os.Args = []string{"uw", "--addr", addr, "--clean"}
		//os.Args = []string{"uw", "--help"}
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
