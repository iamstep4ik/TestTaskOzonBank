package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

type Config struct {
	LogFile    string
	LogLevel   string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	Console    bool
}

func Initialize(cfg Config) error {
	if cfg.MaxSizeMB == 0 {
		cfg.MaxSizeMB = 100
	}

	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		level = zapcore.InfoLevel
	}

	cores := []zapcore.Core{}

	if cfg.LogFile != "" {
		fileCore := zapcore.NewCore(
			getJSONEncoder(),
			getLogWriter(cfg),
			level,
		)
		cores = append(cores, fileCore)
	}

	if cfg.Console || cfg.LogFile == "" {
		consoleCore := zapcore.NewCore(
			getConsoleEncoder(),
			zapcore.AddSync(os.Stderr),
			level,
		)
		cores = append(cores, consoleCore)
	}

	core := zapcore.NewTee(cores...)

	logger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return nil
}

func getJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func getConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func getLogWriter(cfg Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Sync() error {
	return logger.Sync()
}

func GetLogger() *zap.Logger {
	if logger == nil {
		panic("logger not initialized - call log.Initialize() first")
	}
	return logger
}
