package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.ErrorLevel
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)

	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile("log/log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, highPriority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), lowPriority),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	Logger = logger.Sugar()
}
