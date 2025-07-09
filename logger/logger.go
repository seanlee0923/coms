package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func Init() {
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = customTimeEncoder
	consoleEncoderConfig.EncodeLevel = customLevelEncoder
	consoleEncoderConfig.TimeKey = "TIME"
	consoleEncoderConfig.MessageKey = "MESSAGE"
	consoleEncoderConfig.CallerKey = "FILE"
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 로거 코어 설정
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.InfoLevel)

	sugarLogger = zap.New(consoleCore, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

	defer sugarLogger.Sync()
}

var sugarLogger *zap.SugaredLogger

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02-15:04:05.000"))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var levelStr string
	switch level {
	case zapcore.DebugLevel:
		levelStr = "\033[36m[DEBUG]\033[0m" // 청록색
	case zapcore.InfoLevel:
		levelStr = "\033[32m[INFO] \033[0m" // 초록색
	case zapcore.WarnLevel:
		levelStr = "\033[33m[WARN] \033[0m" // 노란색
	case zapcore.ErrorLevel:
		levelStr = "\033[31m[ERROR]\033[0m" // 빨간색
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		levelStr = "\033[35m[FATAL]\033[0m" // 자주색
	}
	enc.AppendString(levelStr)
}

func InfoF(format string, args ...interface{}) {
	sugarLogger.Infof(format, args...)
}

func Info(msg string) {
	sugarLogger.Info(msg)
}

func Error(err error) {
	sugarLogger.Error(err.Error())
}

func Errorf(format string, args ...interface{}) {
	sugarLogger.Errorf(format, args...)
}
