package logic

import (
	"os"
	"time"

	"github.com/richelieu042/chimera/v3/src/file/fileKit"
	"go.uber.org/zap"
)

type (
	cleanJob struct {
		parentDir string
		logger    *zap.Logger
	}
)

func (job *cleanJob) Run() {
	l := job.logger

	l.Warn("[CLEAN] Starts...")
	err := fileKit.Clean(job.parentDir, func(path string, info os.FileInfo) bool {
		birthTime, ok := fileKit.GetBirthTime(info)
		if !ok {
			l.Error("Fallback to ModTime.", zap.String("path", path))
		}

		return time.Since(birthTime) > 10*time.Minute
	})
	if err != nil {
		l.Sugar().Errorf("Clean() failed, error: %+v", err)
	}
	l.Warn("[CLEAN] Ends.")
}
