package main

import (
	"github.com/richelieu-yang/UnchartedWaters/src/cmd"
)

func main() {
	/* test */
	//{
	//	addr := "127.0.0.1:5555"
	//	//addr := "127.0.0.1:5585"
	//	//os.Args = []string{"uw", "--addr", addr, "--clean", "--verbose"}
	//	os.Args = []string{"uw", "--addr", addr, "--clean"}
	//}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
