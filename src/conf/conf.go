package conf

import (
	"image"
	"time"
)

var (
	defSleepInterval = 1_300 * time.Millisecond

	battleSleepInterval = 5_000 * time.Millisecond

	// matchValThreshold 匹配阈值，>=就认为是匹配的
	matchValThreshold float32 = 0.84

	screenshotTimeout = time.Second * 15

	// pointCancelFold 【航海页面】右下角的取消折叠按钮的中心位置
	pointCancelFold = image.Pt(1829, 998)
)

func SetDefSleepInterval(interval time.Duration) {
	defSleepInterval = interval
}

func GetDefSleepInterval() time.Duration {
	return defSleepInterval
}

func GetBattleSleepInterval() time.Duration {
	return battleSleepInterval
}

func GetMatchValThreshold() float32 {
	return matchValThreshold
}

func GetScreenshotTimeout() time.Duration {
	return screenshotTimeout
}

func GetPointCancelFold() image.Point {
	return pointCancelFold
}
