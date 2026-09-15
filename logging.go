package pcommon

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggerContextKey = contextKey("LOGGER")

func CreateLogger() (*zap.Logger, error) {
	var logCfg zap.Config
	if logfmt := os.Getenv("LOG_FORMAT"); logfmt == "dev" {
		logCfg = zap.NewDevelopmentConfig()
	} else {
		logCfg = zap.NewProductionConfig()
	}
	logCfg.OutputPaths = []string{"stdout"}
	logCfg.ErrorOutputPaths = []string{"stdout"}
	return logCfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
}

func MustCreateLogger() *zap.Logger {
	log, err := CreateLogger()
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %s\n", err))
	}
	return log
}

func WithLoggerContextValue(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	value := ctx.Value(loggerContextKey)
	if value == nil {
		panic(errors.New("could not find logger"))
	}

	return value.(*zap.Logger)
}
