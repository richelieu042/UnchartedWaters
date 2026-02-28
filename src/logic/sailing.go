package logic

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/richelieu042/chimera/v3/src/core/pathKit"
	"github.com/richelieu042/chimera/v3/src/core/strKit"
	"github.com/richelieu042/chimera/v3/src/image/imageKit"
	"github.com/richelieu042/chimera/v3/src/ocr/gosseractKit"
	"go.uber.org/zap"
)

// IsSailing 是否在航行页面？
func IsSailing(dirPath, imgPath string, logger *zap.SugaredLogger) (bool, float64, error) {
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
	text, err := gosseractKit.GertText(clippedImgPath)
	if err != nil {
		return false, days, err
	}
	logger.Infof("text: [%s]", text)

	sailing := strKit.Index(text, "航行中") != -1
	logger.Infof("航行中: [%t]", sailing)

	if sailing {
		days, err = getDays(text)
		if err != nil {
			logger.Errorf("获取天数失败, error: %s", err)
		} else {
			logger.Infof("天数: [%.2f]", days)
		}
	}

	return sailing, days, nil
}

// getDays 提取字符串中"天"前面的数字（严格匹配）
func getDays(s string) (float64, error) {
	re := regexp.MustCompile(`(\d+\.?\d*)天`)
	match := re.FindStringSubmatch(s)

	if len(match) > 1 {
		return strconv.ParseFloat(match[1], 64)
	}
	return 0, fmt.Errorf("invalid string: %s", s)
}
