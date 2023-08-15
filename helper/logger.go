package helper

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

var logger *zap.Logger

func init() {

	var w = zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/log.txt",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
	})

	pe := zap.NewProductionEncoderConfig()
	pe.TimeKey = "timestamp"
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	jsonEncoder := zapcore.NewJSONEncoder(pe)

	fileLogger := zapcore.NewCore(jsonEncoder, w, zap.DebugLevel)
	consoleLogger := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

	loggerCores := []zapcore.Core{consoleLogger}
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		loggerCores = append(loggerCores, fileLogger)
	}

	core := zapcore.NewTee(loggerCores...)

	logger = zap.New(core)
}

func WaitForLoggerComplete() {
	logger.Sync()
}

func GetLogger(ctx *gin.Context) *Logger {
	return ctx.MustGet("logger").(*Logger)
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
	_, file, line, _ := runtime.Caller(2)

	return l.logger.With(zap.String("requestID", l.uuid)).
		With(zap.String("file", file)).
		With(zap.Int("line", line))
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

func (l *Logger) With(key string, value interface{}) {
	l.logger = l.logger.With(key, value)
}

func (l *Logger) Close() {
	l.logger.Sync()
}
