package main

import (
	"os"

	"github.com/richelieu-yang/UnchartedWaters/src/cmd"
)

func main() {
	os.Args = []string{"uw", "-h"}
	//os.Args = []string{"uw", "clean"}
	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	//ins, err := adbKit.NewInstance("127.0.0.1:5585", true, true)
	//if err != nil {
	//	console.Fatalf("adbKit.NewInstance() failed: %s", err)
	//}
	//w, h, err := ins.GetPhysicalSize()
	//if err != nil {
	//	console.Fatalf("ins.GetPhysicalSize() failed: %s", err)
	//}
	//console.Infof("physical size: %dx%d", w, h)
	//if w != 1920 && h != 1080 {
	//	console.Fatalf("unsupported physical size: %dx%d", w, h)
	//}

	//a := adb.New("127.0.0.1:5555", "1920x1080")
	//if err := a.CheckEnv(); err != nil {
	//	console.Fatalf("a.CheckEnv() failed: %s", err)
	//	return
	//}
	//if err := a.Initialize(); err != nil {
	//	console.Fatalf("a.Initialize() failed: %s", err)
	//	return
	//}

	//for {
	//	// 随机睡眠500ms ~ 1000ms
	//	randInt := randomKit.Int(1000, 2001)
	//	time.Sleep(time.Millisecond * time.Duration(randInt))
	//
	//	imgName := "screenshot_" + timeKit.FormatCurrent(timeKit.FormatFileName) + ".png"
	//	start := time.Now()
	//	if err := ins.Screenshot(imgName); err != nil {
	//		console.Errorf("a.Screenshot() failed: %s", err)
	//		continue
	//	}
	//	console.Infof("screenshot: %s, cost: %s", imgName, time.Since(start))
	//}
}
