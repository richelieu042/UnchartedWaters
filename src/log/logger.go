package log

import (
	"os"

	"github.com/richelieu042/chimera/v3/src/log/zapKit"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	level = zap.DebugLevel
)

func SetLevel(l zapcore.Level) {
	level = l
}

func NewLogger(prefix string) *zap.Logger {
	enc := zapKit.NewEncoder(zapKit.WithEncoderMessagePrefix(prefix))
	core := zapKit.NewCore(enc, os.Stdout, level)
	return zapKit.NewLogger(core)
}
