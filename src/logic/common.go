package logic

import (
	"context"
	"image"
	"time"

	"github.com/richelieu-yang/UnchartedWaters/src/conf"
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/concurrency/rateLimitKit"
	"github.com/richelieu042/chimera/v3/src/gocvKit"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

var (
	// opLimiter 限流器：避免短时间内操作太多次
	opLimiter = rateLimitKit.NewUberLimiter(1, ratelimit.Per(800*time.Millisecond), ratelimit.WithoutSlack)
)

func match(imgPath, templPath string) (matchVal float32, matchRect image.Rectangle, err error) {
	return gocvKit.MatchTemplate(imgPath, templPath, gocv.TmCcoeffNormed, true)
	//return gocvKit.MatchTemplate(imgPath, templPath, gocv.TmCcoeffNormed)
}

/*
@return

	0: 匹配失败
	1: 匹配成功，匹配度低
	2: 匹配成功，匹配度高，点击失败
	3: 匹配成功，匹配度高，点击成功
*/
func matchAndTap(adbClient adbKit.Client, l *zap.Logger, op, imgPath, templPath string) int {
	matchVal, matchRect, err := match(imgPath, templPath)
	if err != nil {
		l.Sugar().Errorf("Fail to match template, op: %s, error: %+v", op, err)
		return 0
	}

	if matchVal < conf.GetMatchValThreshold() {
		l.Debug("MatchVal is too low.", zap.String("op", op), zap.Float32("matchVal", matchVal))
		return 1
	}
	l.Debug("MatchVal is enough.", zap.String("op", op), zap.Float32("matchVal", matchVal))

	x := matchRect.Min.X + matchRect.Dx()/2
	y := matchRect.Min.Y + matchRect.Dy()/2
	if err := _tap(adbClient, x, y, 6); err != nil {
		l.Error("Fail to tap.", zap.String("op", op), zap.Int("x", x), zap.Int("y", y), zap.Error(err))
		return 2
	}
	l.Info("Manager to tap.", zap.String("op", op), zap.Int("x", x), zap.Int("y", y))
	return 3
}

func _tap(adbClient adbKit.Client, x, y int, axisOffset int) error {
	opLimiter.Take() // 先等一会
	return adbClient.TapAsHumanBeings(context.TODO(), x, y, axisOffset)
}
