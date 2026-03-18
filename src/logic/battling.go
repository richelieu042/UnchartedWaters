package logic

import (
	"github.com/richelieu042/chimera/v3/src/core/error/errKit"
	"go.uber.org/zap"
)

// isBattling 是否在海战页面？
func isBattling(logger *zap.Logger, imgPath string) (bool, error) {
	templPath := "images/battle/auto.png"
	matchVal, _, err := match(imgPath, templPath)
	if err != nil {
		return false, errKit.Wrapf(err, "fail to math with template(path: %s)", templPath)
	}
	if matchVal >= 0.85 {
		return true, nil
	}
	logger.Debug("MatchVal is too low.", zap.Float32("matchVal", matchVal), zap.String("templPath", templPath))

	templPath = "images/battle/not_auto.png"
	matchVal, _, err = match(imgPath, templPath)
	if err != nil {
		return false, errKit.Wrapf(err, "fail to math with template(path: %s)", templPath)
	}
	if matchVal >= 0.85 {
		return true, nil
	}
	logger.Debug("MatchVal is too low.", zap.Float32("matchVal", matchVal), zap.String("templPath", templPath))

	return false, nil
}
