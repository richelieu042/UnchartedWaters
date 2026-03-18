package logic

import (
	"time"

	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/concurrency/rateLimitKit"
	"github.com/richelieu042/chimera/v3/src/core/error/errKit"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

// isBattling 是否在海战页面？
func isBattling(logger *zap.Logger, imgPath string) (bool, error) {
	templPath := "images/battle/flag.png"
	matchVal, _, err := match(imgPath, templPath)
	if err != nil {
		return false, errKit.Wrapf(err, "fail to math with template(path: %s)", templPath)
	}
	if matchVal < 0.85 {
		logger.Debug("MatchVal is too low.", zap.Float32("matchVal", matchVal), zap.String("templPath", templPath))
		return false, nil
	}

	return true, nil
}

func processBattling(adbClient adbKit.Client, l *zap.Logger, imgPath string) {
	// 限流器：避免短时间内操作太多次
	limiter := rateLimitKit.NewUberLimiter(1, ratelimit.Per(500*time.Millisecond), ratelimit.WithoutSlack)

	// TODO: (0) 开启“自动战斗”
	{
		op := "enable_auto"
		templPath := "images/sail/not_auto.png.png"

		flag := matchAndTap(adbClient, l, op, imgPath, templPath, limiter)
		switch flag {
		case 0, 2:
			return // 中断
		case 1, 3:
			// 继续向下走
		default:
			l.Panic("Unknown flag.", zap.String("op", op), zap.Int("flag", flag))
		}
	}

	// (1) 召唤海军
	{
		op := "navy"
		templPath := "images/sail/navy.png"

		flag := matchAndTap(adbClient, l, op, imgPath, templPath, limiter)
		switch flag {
		case 0, 2:
			return // 中断
		case 1, 3:
			// 继续向下走
		default:
			l.Panic("Unknown flag.", zap.String("op", op), zap.Int("flag", flag))
		}
	}

	// (2) 召唤副舰
	{
		op := "assistant"
		templPath := "images/sail/assistant.png"

		flag := matchAndTap(adbClient, l, op, imgPath, templPath, limiter)
		switch flag {
		case 0, 2:
			return // 中断
		case 1, 3:
			// 继续向下走
		default:
			l.Panic("Unknown flag.", zap.String("op", op), zap.Int("flag", flag))
		}
	}
}
