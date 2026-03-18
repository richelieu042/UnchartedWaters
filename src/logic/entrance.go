package logic

import (
	"time"

	"github.com/richelieu-yang/UnchartedWaters/src/log"
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/core/pathKit"
	"github.com/richelieu042/chimera/v3/src/core/strKit"
	"github.com/richelieu042/chimera/v3/src/file/fileKit"
	"github.com/richelieu042/chimera/v3/src/randomKit"
	"github.com/richelieu042/chimera/v3/src/time/timeKit"
	"go.uber.org/zap"
)

const (
	defInterval = 1_000
)

var (
	// 单位：ms
	sleepInterval = defInterval
)

func Start(adbClient adbKit.Client, logger *zap.Logger) {
	// 目前仅支持 1920x1080 尺寸
	w, h, err := adbClient.GetPhysicalSize()
	if err != nil {
		logger.Sugar().Panicf("Fail to get physical size, error: %+v", err)
		return
	}
	logger.Info("physical size", zap.Int("width", w), zap.Int("height", h))
	if w != 1920 || h != 1080 {
		logger.Panic("Size is unsupported!!!", zap.Int("width", w), zap.Int("height", h))
		return
	}

	for {
		// 随机睡眠
		randInt := randomKit.Int(-200, 201) + sleepInterval
		duration := time.Millisecond * time.Duration(randInt)
		logger.Info("Sleep starts...................................................", zap.String("duration", duration.String()))
		time.Sleep(duration)
		logger.Info("Sleep ends.", zap.String("duration", duration.String()))

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
		logger.Sugar().Infof("dirPath: [%s]", dirPath)

		imgPath := pathKit.Join(dirPath, "screenshot.png")
		if err := adbClient.Screenshot(imgPath); err != nil {
			logger.Error("Screenshot() failed", zap.Error(err))
			continue
		}

		/* （1）航行中 */
		{
			l := log.NewLogger("[SAILING] ")
			l = l.With(zap.String("dirPath", dirPath))

			flag, days, err := isSailing(l, dirPath, imgPath)
			if err != nil {
				l.Sugar().Errorf("isSailing() fails, error: %+v", err)
				continue
			} else {
				l.Sugar().Infof("Is sailing? [%t]", flag)
				if flag {
					processSailing(adbClient, l, imgPath, days)
					continue
				}
			}
		}

		/* （2）战斗中 */
		{
			l := log.NewLogger("[BATTLING] ")
			l = l.With(zap.String("dirPath", dirPath))

			flag, err := isBattling(l, imgPath)
			if err != nil {
				l.Sugar().Errorf("isBattling() fails, error: %+v", err)
				continue
			} else {
				l.Sugar().Infof("Is battling? [%t]", flag)
				if flag {
					processBattling(adbClient, l, imgPath)
					continue
				}
			}
		}
	}
}
