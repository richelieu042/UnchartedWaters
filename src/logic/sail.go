package logic

import (
	"fmt"
	"image"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/core/pathKit"
	"github.com/richelieu042/chimera/v3/src/core/strKit"
	"github.com/richelieu042/chimera/v3/src/image/imageKit"
	"github.com/richelieu042/chimera/v3/src/ocr/gosseractKit"
	"github.com/richelieu042/chimera/v3/src/randomKit"
	"go.uber.org/zap"
)

// isSailing 是否在航行页面？
func isSailing(logger *zap.Logger, dirPath, imgPath string) (bool, float64, error) {
	x0 := 748
	y0 := 995
	x1 := 962
	y1 := 1040

	days := -1.0 // -1.0: 获取天数失败

	clippedImgPath := pathKit.Join(dirPath, "clipped.png")
	err := imageKit.ClipWithPath(imgPath, clippedImgPath, x0, y0, x1-x0+1, y1-y0+1)
	if err != nil {
		return false, days, err
	}
	text, err := gosseractKit.GertText(clippedImgPath, "chi_sim")
	if err != nil {
		return false, days, err
	}
	logger.Sugar().Debugf("text: [%s]", text)

	//sailing := strKit.Index(text, "航行中") != -1
	//if sailing {
	//	days, err = getDays(text)
	//	if err != nil {
	//		logger.Error("获取剩余航行天数失败！", zap.Error(err))
	//	} else {
	//		logger.Sugar().Debugf("天数: [%.2f]", days)
	//	}
	//}

	sailing := false
	days, err = getDays(text)
	if err != nil {
		logger.Error("获取剩余航行天数失败！", zap.Error(err))
	} else {
		logger.Sugar().Debugf("天数: [%.2f]", days)
		sailing = days > 0
	}

	return sailing, days, nil
}

// getDays 提取字符串中"天"前面的数字（严格匹配）
func getDays(s string) (float64, error) {
	if err := strKit.AssertNotBlank(s, "s"); err != nil {
		return -1, err
	}

	re := regexp.MustCompile(`(\d+\.?\d*)天`)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		return strconv.ParseFloat(match[1], 64)
	}
	return 0, fmt.Errorf("invalid string: %s", s)
}

// processSailing 处理航行中
/*
@param days 剩余航行天数
*/
func processSailing(adbClient adbKit.Client, l *zap.Logger, imgPath string, days float64) {
	if days < 0 {
		l.Warn("Fail to get left days.", zap.Float64("days", days))
		return
	} else if days < 0.2 {
		l.Info("Left days is too few, do nothing.", zap.Float64("days", days))
		return
	}

	l.Info("Left days is enough.", zap.Float64("days", days))

	// (1) 【高优先级】送礼
	{
		op := "gift"
		templPath := "images/sail/gift.png"

		flag := matchAndTap(adbClient, l, op, imgPath, templPath)
		switch flag {
		case 0, 2:
			return
		case 1, 3: // 继续向下走
		default:
			l.Panic("Unknown flag.", zap.String("op", op), zap.Int("flag", flag))
		}
	}

	// (2) 模拟点击 3 处高频事件点
	points := []image.Point{
		{X: 1168, Y: 708},
		{X: 1168, Y: 861},
		{X: 1315, Y: 706},
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
				l.Error("Fail to tap.", zap.String("op", "event"), zap.Int("index", i))
				return
			}

			l.Info("Manager to tap.", zap.String("op", "event"), zap.Int("index", i))
		}()
	}
	wg.Wait()

	// (3) 测量
	{
		op := "measure"
		templPath := "images/sail/measure.png"

		flag := matchAndTap(adbClient, l, op, imgPath, templPath)
		switch flag {
		case 0, 2:
			return
		case 1:
			// 继续向下走
		case 3:
			return // 此种情况下，不需要走下一步
		default:
			l.Panic("Unknown flag.", zap.String("op", op), zap.Int("flag", flag))
		}
	}

	// (4) 取消（右下角的）折叠
	// TODO: 有时会误触，可能原因：折叠和取消折叠的按钮太像了
	{
		op := "cancel_fold"
		templPath := "images/sail/folded.png"

		flag := matchAndTap(adbClient, l, op, imgPath, templPath)
		switch flag {
		case 0, 2:
			return
		case 1:
			// 未折叠，继续向下走
		case 3:
			return // 折叠已取消，等下个循环
		default:
			l.Panic("Unknown flag.", zap.String("op", op), zap.Int("flag", flag))
		}
	}
}
