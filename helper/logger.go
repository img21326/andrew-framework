package helper

import (
	"context"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

var logger *zap.Logger

func init() {

	var w = zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log.txt",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
	})

	pe := zap.NewProductionEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.DebugLevel,
		),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	logger = zap.New(core)
}

func WaitForLoggerComplete() {
	logger.Sync()
}

func GetLogger() *zap.SugaredLogger {
	return logger.Sugar()
}

type Logger struct {
	uuid   string
	logger *zap.SugaredLogger
}

func NewLogger(uuid string) *Logger {
	return &Logger{
		uuid:   uuid,
		logger: logger.Sugar(),
	}
}

func (l *Logger) log(_ context.Context) *zap.SugaredLogger {
	return l.logger.With(zap.String("requestID", l.uuid))
}

func (l *Logger) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log(ctx).Infof(s, args...)
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log(ctx).Warnf(s, args...)
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.log(ctx).Errorf(s, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	l.log(ctx).Debugf("%s [%s]", sql, elapsed)
}
