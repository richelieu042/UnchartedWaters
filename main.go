package main

import (
	"time"

	"github.com/richelieu-yang/UnchartedWaters/src/adb"
	"github.com/richelieu-yang/chimera/v3/src/log/console"
	"github.com/richelieu-yang/chimera/v3/src/randomKit"
	"github.com/richelieu-yang/chimera/v3/src/time/timeKit"
)

func main() {
	a := adb.New("127.0.0.1:5555", "1920x1080")
	if err := a.CheckEnv(); err != nil {
		console.Fatalf("a.CheckEnv() failed: %s", err)
		return
	}
	if err := a.Initialize(); err != nil {
		console.Fatalf("a.Initialize() failed: %s", err)
		return
	}

	for {
		// 随机睡眠500ms ~ 1000ms
		randInt := randomKit.Int(500, 1001)
		time.Sleep(time.Millisecond * time.Duration(randInt))

		imgName := "screenshot_" + timeKit.FormatCurrent(timeKit.FormatFileName) + ".png"
		start := time.Now()
		if err := a.Screenshot(imgName); err != nil {
			console.Errorf("a.Screenshot() failed: %s", err)
			continue
		}
		console.Infof("screenshot: %s, cost: %s", imgName, time.Since(start))

		//// adb -s {a.address} exec-out screencap -p
		//data, err := cmdKit.RunCombinedly(context.TODO(), "adb", "-s", a.address, "exec-out", "screencap", "-p")
		//if err != nil {
		//	console.Errorf("cmdKit.RunCombinedly() failed: %s", err)
		//	continue
		//}
		//
		//err = os.WriteFile(imgName, data, 0644)
		//if err != nil {
		//	panic(err)
		//}
	}
}
