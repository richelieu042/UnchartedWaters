package logic

import (
	"image"
	"sync"
	"time"

	"github.com/richelieu-yang/UnchartedWaters/src/conf"
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/core/error/errKit"
	"github.com/richelieu042/chimera/v3/src/randomKit"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// isBattling 是否在海战页面？
func isBattling(logger *zap.Logger, imgPath string) (bool, error) {
	templPath := "images/battle/flag.png"
	op := "isBattling"

	matchVal, _, err := match(imgPath, templPath)
	if err != nil {
		return false, errKit.Wrapf(err, "fail to math with template(path: %s)", templPath)
	}
	logger.Sugar().Debugf("matchVal: [%.2f], op: [%s]", matchVal, op)
	if matchVal < conf.GetMatchValThreshold() {
		logger.Debug("MatchVal is too low.", zap.Float32("matchVal", matchVal), zap.String("templPath", templPath))
		return false, nil
	}

	return true, nil
}

func processBattling(adbClient adbKit.Client, l *zap.Logger, imgPath string, tapCount *atomic.Int32) {
	// (1) 开启“自动战斗”
	{
		op := "enable_auto"
		templPath := "images/battle/not_auto.png"

		flag := matchAndTap(adbClient, l, op, imgPath, templPath)
		switch flag {
		case 0, 2:
			return // 中断
		case 1:
			l.Info("Already auto.")
			// 继续向下走
		case 3:
			// 继续向下走
			l.Info("Manager to enable auto.")
		default:
			l.Panic("Unknown flag.", zap.String("op", op), zap.Int("flag", flag))
		}
	}

	/* (2) 召唤友舰（海军可以在设置中自动召唤） */
	if tapCount.Load() > 0 {
		tapCount.Dec()
		l.Info("tapCount - 1", zap.Int32("left_count", tapCount.Load()))

		points := []image.Point{
			{1818, 333},
			//{1708, 333},
		}
		var wg sync.WaitGroup
		for i, p := range points {
			wg.Add(1)

			go func() {
				defer wg.Done()

				// 随机等一会，使顺序更加随机
				ri := randomKit.Int(10, 21)
				time.Sleep(time.Millisecond * time.Duration(ri))

				if err := _tap(adbClient, p.X, p.Y, 10); err != nil {
					l.Error("Fail to tap.", zap.Int("i", i))
					return
				}
				l.Info("Manager to tap.", zap.Int("i", i))
			}()
		}
		wg.Wait()
	} else {
		l.Warn("Count is zero!")
	}
}
