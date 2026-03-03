package main

import (
	"os"

	"github.com/richelieu-yang/UnchartedWaters/src/cmd"
	"github.com/richelieu042/chimera/v3/src/log/console"
	"go.uber.org/zap"
)

func main() {
	/* test */
	{
		addr := "127.0.0.1:5585"
		os.Args = []string{"uw", "--addr", addr, "--clean", "--verbose"}
	}

	if err := cmd.Execute(); err != nil {
		console.Panic("Execute failed", zap.Error(err))
	}
}
