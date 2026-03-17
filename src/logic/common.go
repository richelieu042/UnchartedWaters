package logic

import (
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/gocvKit"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

/*
@return

	0: 匹配失败
	1: 匹配成功，匹配度低
	2: 匹配成功，匹配度高，点击失败
	3: 匹配成功，匹配度高，点击成功
*/
func matchAndTap(adbClient adbKit.Client, l *zap.Logger, op, imgPath, templPath string, limiter ratelimit.Limiter) int {
	matchVal, matchRect, err := gocvKit.MatchTemplate(imgPath, templPath, gocv.TmCcoeffNormed)
	if err != nil {
		l.Sugar().Errorf("Fail to match template, op: %s, error: %+v", op, err)
		return 0
	}
	l.Sugar().Infof("op: %s, matchVal: %.2f", op, matchVal)

	if matchVal < 0.8 {
		return 1
	}

	limiter.Take() // 等一会
	if err := adbClient.TapAsHumanBeings(matchRect.Min.X+matchRect.Dx()/2, matchRect.Min.Y+matchRect.Dy()/2, 10); err != nil {
		l.Error("Fail to tap as human beings.", zap.String("op", op))
		return 2
	}
	l.Info("Manager to tap as human beings.", zap.String("op", op))
	return 3
}
