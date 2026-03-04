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

func Start(adbClient adbKit.Client, logger *zap.Logger) {
	// 目前仅支持 1920x1080 尺寸
	w, h, err := adbClient.GetPhysicalSize()
	if err != nil {
		logger.Panic("Fail to get physical size.", zap.Error(err))
		return
	}
	logger.Info("physical size", zap.Int("width", w), zap.Int("height", h))
	if w != 1920 || h != 1080 {
		logger.Panic("unsupported size.", zap.Int("width", w), zap.Int("height", h))
		return
	}

	for {
		// 随机睡眠
		randInt := randomKit.Int(1000, 2001)
		duration := time.Millisecond * time.Duration(randInt)
		console.Info("Sleep starts...................................................", zap.String("duration", duration.String()))
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
			timeKit.Format(now, "2006-01-02"),
			timeKit.Format(now, "15.04"),
			timeKit.Format(now, "05.000"),
		)
		if err := fileKit.MkDirs(dirPath); err != nil {
			logger.Error("MkDirs() failed", zap.Error(err))
			continue
		}
		console.Infof("dirPath: [%s]", dirPath)

		// 子logger，后续的输出都用它
		l := logger.With(zap.String("dirPath", dirPath))

		imgPath := pathKit.Join(dirPath, "screenshot.png")
		if err := adbClient.Screenshot(imgPath); err != nil {
			l.Error("Screenshot() failed", zap.Error(err))
			continue
		}

		{
			sailFlag, days, err := IsSailing(dirPath, imgPath, l)
			if err != nil {
				l.Error("IsSailing fails.", zap.Error(err))
			} else {
				if sailFlag {
					if days < 0 {
						l.Warn("Is sailing, but fail to get left days.", zap.Float64("days", days))
					} else if days < 0.5 {
						l.Info("Is sailing, but left days is too few, do nothing.", zap.Float64("days", days))
					} else {
						// TODO:
						l.Info("Is sailing.", zap.Float64("days", days))

						// 模拟点击
						points := []*adbKit.Point{
							{X: 1168, Y: 708},
							{X: 1168, Y: 861},
							{X: 1315, Y: 706},
						}
						for i, point := range points {
							if err := adbClient.TapAsHumanBeings(point.X, point.Y, 10); err != nil {
								l.Error("Fail to tap as human beings.", zap.Int("index", i))
							} else {
								l.Info("Manager to tap as human beings.", zap.Int("index", i))
							}
							time.Sleep(time.Millisecond * 500)
						}
					}
					continue
				}
				console.Info("Not sailing.")
			}
		}

		//start := time.Now()
		//if err := ins.Screenshot(imgName); err != nil {
		//	console.Errorf("a.Screenshot() failed: %s", err)
		//	continue
		//}
		//console.Infof("screenshot: %s, cost: %s", imgName, time.Since(start))
	}
}
