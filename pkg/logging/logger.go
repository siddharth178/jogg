package logging

import (
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var configToLevel = map[string]zapcore.Level{
	"DEBUG":  zapcore.DebugLevel,
	"INFO":   zapcore.InfoLevel,
	"WARN":   zapcore.WarnLevel,
	"ERROR":  zapcore.ErrorLevel,
	"DPANIC": zapcore.DPanicLevel,
	"PANIC":  zapcore.PanicLevel,
	"FATAL":  zapcore.FatalLevel,
}

const LoggerDefaultLevel = zapcore.InfoLevel

func getLogLevel(level string) (zapcore.Level, error) {
	level = strings.ToUpper(level)
	if logLevel, ok := configToLevel[level]; ok {
		return logLevel, nil
	}
	return LoggerDefaultLevel, errors.New("invalid log level")
}

func GetLogger(cfg config.LogConfig) (*zap.SugaredLogger, error) {
	level, err := getLogLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	loggerCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableStacktrace: true, // don't use logger's facility to capture stack
		Encoding:          cfg.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return logger.Sugar(), nil
}
