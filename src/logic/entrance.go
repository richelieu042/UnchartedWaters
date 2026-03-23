package logic

import (
	"time"

	"github.com/richelieu-yang/UnchartedWaters/src/log"
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"github.com/richelieu042/chimera/v3/src/atomic/atomicKit"
	"github.com/richelieu042/chimera/v3/src/core/pathKit"
	"github.com/richelieu042/chimera/v3/src/core/strKit"
	"github.com/richelieu042/chimera/v3/src/file/fileKit"
	"github.com/richelieu042/chimera/v3/src/randomKit"
	"github.com/richelieu042/chimera/v3/src/time/timeKit"
	"go.uber.org/zap"
)

const (
	defSleepInterval    = 800
	battleSleepInterval = 5_000

	// 默认的次数（在战斗中点击按钮的次数）
	defBattleTapCount = 5
)

func Start(adbClient adbKit.Client, logger *zap.Logger, disableSail, disableBattle bool) {
	sleepInterval := defSleepInterval                       // 单位：ms
	battleTapCount := atomicKit.NewInt32(defBattleTapCount) // 第1次点击大概率无效，因为刚开战按钮在CD
	preSailing := false                                     // 上一次是否是在航行？

	/* 尺寸，目前仅支持 1920x1080 */
	w, h, err := adbClient.GetPhysicalSize()
	if err != nil {
		logger.Sugar().Panicf("Fail to get physical size, error: %+v", err)
		return
	}
	logger.Info("physical size", zap.Int("width", w), zap.Int("height", h))
	if w < h {
		w, h = h, w
	}
	if w != 1920 || h != 1080 {
		logger.Panic("Size isn't supported!!!", zap.Int("width", w), zap.Int("height", h))
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
		start := time.Now()
		if err := adbClient.Screenshot(imgPath); err != nil {
			logger.Error("Screenshot() failed", zap.Error(err))
			continue
		}
		logger.Debug("Screenshot succeeds.", zap.String("duration", time.Since(start).String()))

		/* （1）航行中 */
		{
			l := log.NewLogger("[SAILING] ")

			flag, days, err := isSailing(l, dirPath, imgPath)
			if err != nil {
				l.Sugar().Errorf("isSailing() fails, error: %+v", err)
				continue
			} else {
				l.Sugar().Infof("Is sailing? [%t]", flag)
				if flag {
					sleepInterval = defSleepInterval
					battleTapCount.Store(defBattleTapCount) // 重置战斗次数，因为已经脱离了战斗
					preSailing = true

					if disableSail {
						l.Info("Sail is disabled.")
					} else {
						processSailing(adbClient, l, imgPath, days)
						continue
					}
				}
			}
		}

		/* （2）战斗中 */
		{
			l := log.NewLogger("[BATTLING] ")

			flag, err := isBattling(l, imgPath)
			if err != nil {
				l.Sugar().Errorf("isBattling() fails, error: %+v", err)
				continue
			} else {
				l.Sugar().Infof("Is battling? [%t]", flag)
				if flag {
					sleepInterval = battleSleepInterval // 时间长一点，以避免无效点击
					if preSailing {
						preSailing = false
						d := time.Second * 8
						l.Warn("初次进入战斗，sleep starts...", zap.String("duration", d.String()))
						time.Sleep(d) // 刚进入战斗，前8s内不允许召唤海军or友舰
						l.Warn("初次进入战斗，sleep ends.")
						continue
					}

					if disableBattle {
						l.Info("Battle is disabled.")
					} else {
						processBattling(adbClient, l, imgPath, battleTapCount)
						continue
					}
				}
			}
		}
	}
}
