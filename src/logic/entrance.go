package logic

import (
	"time"

	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/core/pathKit"
	"github.com/richelieu042/chimera/v3/src/core/strKit"
	"github.com/richelieu042/chimera/v3/src/file/fileKit"
	"github.com/richelieu042/chimera/v3/src/log/console"
	"github.com/richelieu042/chimera/v3/src/randomKit"
	"github.com/richelieu042/chimera/v3/src/time/timeKit"
	"go.uber.org/zap"
)

func Start(adbClient adbKit.Client, logger *zap.SugaredLogger) {
	// 目前仅支持 1920x1080 尺寸
	w, h, err := adbClient.GetPhysicalSize()
	if err != nil {
		logger.Panic("Fail to get physical size.", zap.Error(err))
		return
	}
	logger.Infof("physical size: [%dx%d]", w, h)
	if w != 1920 || h != 1080 {
		logger.Panic("unsupported size.", zap.Int("width", w), zap.Int("height", h))
		return
	}

	for {
		var l *zap.Logger

		l.With()

		logger.With()

		// 随机睡眠
		randInt := randomKit.Int(1000, 2001)
		duration := time.Millisecond * time.Duration(randInt)
		console.Info("Sleep starts................................", zap.String("duration", duration.String()))
		time.Sleep(duration)
		console.Info("Sleep ends.", zap.String("duration", duration.String()))

		now := time.Now()
		/*
			第1层：__tmp
			第2层：adb连接地址（处理了":"）
			第3层：时间精确到“分钟”
			第4层：时间精确到“毫秒”
		*/
		dirPath := pathKit.Join("__tmp",
			strKit.ReplaceAll(adbClient.GetAddress(), ":", "_"),
			timeKit.Format(now, "2006-01-02T15.04"),
			timeKit.Format(now, "05.000"),
		)
		if err := fileKit.MkDirs(dirPath); err != nil {
			logger.Errorf("fileKit.MkDirs() failed: %s", err)
			continue
		}
		console.Infof("dirPath: [%s]", dirPath)

		imgPath := pathKit.Join(dirPath, "screenshot.png")
		if err := adbClient.Screenshot(imgPath); err != nil {
			logger.Errorf("adbClient.Screenshot() failed: %s", err)
			continue
		}

		{
			_, _, _ = IsSailing(dirPath, imgPath, logger)
		}

		//start := time.Now()
		//if err := ins.Screenshot(imgName); err != nil {
		//	console.Errorf("a.Screenshot() failed: %s", err)
		//	continue
		//}
		//console.Infof("screenshot: %s, cost: %s", imgName, time.Since(start))
	}
}
