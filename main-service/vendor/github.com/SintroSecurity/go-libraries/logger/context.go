package logger

import (
	"context"
	"errors"
)

type ctxKeyLoggerType uint

const ctxKeyLogger ctxKeyLoggerType = iota + 1

func GetLoggerFromContext(ctx context.Context) *Logger {
	maybeLogger, ok := ctx.Value(ctxKeyLogger).(*Logger)
	if !ok {
		panic(errors.New("no logger in context"))
	}
	loggerCopy := *maybeLogger
	return &loggerCopy
}

func NewContextWithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}
