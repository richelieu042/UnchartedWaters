package logic

import (
	"time"

	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/log/console"
	"github.com/richelieu042/chimera/v3/src/randomKit"
	"go.uber.org/zap"
)

func Start(adbClient adbKit.Client, logger *zap.SugaredLogger) {
	// 目前仅支持 1920x1080 尺寸
	w, h, err := adbClient.GetPhysicalSize()
	if err != nil {
		panic(err)
	}
	logger.Infof("physical size: [%dx%d]", w, h)
	if w != 1920 || h != 1080 {
		logger.Panicf("unsupported size: %dx%d", w, h)
	}

	for {
		// 随机睡眠
		randInt := randomKit.Int(800, 1301)
		duration := time.Millisecond * time.Duration(randInt)
		console.Info("Sleep starts................................", zap.String("duration", duration.String()))
		time.Sleep(duration)
		console.Info("Sleep ends.", zap.String("duration", duration.String()))

		//imgName := "screenshot_" + timeKit.FormatCurrent(timeKit.FormatFileName) + ".png"
		//start := time.Now()
		//if err := ins.Screenshot(imgName); err != nil {
		//	console.Errorf("a.Screenshot() failed: %s", err)
		//	continue
		//}
		//console.Infof("screenshot: %s, cost: %s", imgName, time.Since(start))
	}
}
