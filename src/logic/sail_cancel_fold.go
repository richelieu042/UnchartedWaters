package logic

import (
	"context"
	"sync"
	"time"

	"github.com/richelieu-yang/UnchartedWaters/src/conf"
	"github.com/richelieu042/chimera/v3/src/android/adbKit"
	"go.uber.org/zap"
)

type result struct {
	Found     bool
	MatchVal  float32
	TemplPath string
}

// cancelFoldWhenSailing
/*
	逻辑：
		匹配三张图片，如果找到的数量>=2，认为此时右下角的菜单未被折叠，结束；
		否则，点击右下角的折叠按钮.
*/
func cancelFoldWhenSailing(adbClient adbKit.Client, l *zap.Logger, imgPath string) {
	op := "cancel_fold"
	templPaths := []string{
		"images/sail/guancha.png",
		"images/sail/renshi.png",
		"images/sail/wangfeng.png",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var wg sync.WaitGroup
	ch := make(chan *result, len(templPaths))

	for i, templPath := range templPaths {
		wg.Add(1)

		go worker(ctx, l, &wg, op, ch, imgPath, i, templPath)
	}

	// 所有 worker 结束后关闭 channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	tapFlag := false
	received := 0
	notFoundCount := 0
loop:
	for {
		select {
		case <-ctx.Done():
			// 条件 1：超过 3s（context timeout）
			l.Warn("Timeout.", zap.String("op", op))
			break loop

		case res, ok := <-ch:
			if !ok {
				// 条件 2：channel 关闭 = 3 个 goroutine 都结束且没人发现结果
				l.Warn("All workers end.", zap.String("op", op))
				break loop
			}

			received++
			if !res.Found {
				notFoundCount++
				// 条件3：有2个图标不存在，认为被隐藏了，需要点击"取消隐藏"按钮
				if notFoundCount >= 2 {
					tapFlag = true
					cancel() // 通知其他 worker 退出
					break loop
				}
			}
			if received == len(templPaths) {
				// 条件4：所有结果都收到了，但不满足tap的条件
				l.Info("Receive all results, needn't to tap.")
				break loop
			}
		}
	}

	l.Sugar().Infof("tapFlag: %t, notFoundCount: %d", tapFlag, notFoundCount)
	if !tapFlag {
		return
	}
	p := conf.GetPointCancelFold()
	err := _tap(adbClient, p.X, p.Y, 6)
	if err != nil {
		l.Error("Fail to tap.", zap.String("op", op), zap.Error(err))
		return
	}
	l.Info("Manager to tap.", zap.String("op", op))
}

func worker(ctx context.Context, l *zap.Logger, wg *sync.WaitGroup, op string, ch chan<- *result, imgPath string, i int, templPath string) {
	defer wg.Done()

	matchVal, _, err := match(imgPath, templPath)
	if err != nil {
		l.Error("Fail to match.", zap.String("op", op), zap.Int("index", i), zap.Error(err))
		return
	}
	rst := &result{
		Found:     matchVal >= conf.GetMatchValThreshold(),
		MatchVal:  matchVal,
		TemplPath: templPath,
	}
	l.Info("Manager to match.", zap.String("op", op), zap.Int("index", i), zap.Bool("found", rst.Found), zap.Float32("matchVal", rst.MatchVal))

	// 先排除 ctx 已经取消的情况
	if ctx.Err() != nil {
		l.Warn("Context is done, result ignored.", zap.String("op", op), zap.Int("index", i))
		return
	}
	// 再进入 select 等待（小瑕疵：仍然有那个极小窗口——两个 select 之间 ctx 恰好取消，第二个 select 还是可能随机走发送分支）
	select {
	case <-ctx.Done():
		l.Warn("Context is done, result ignored.", zap.String("op", op), zap.Int("index", i))
		return
	case ch <- rst:
		l.Info("Manager to send result.", zap.String("op", op), zap.Int("index", i))
	}
}
