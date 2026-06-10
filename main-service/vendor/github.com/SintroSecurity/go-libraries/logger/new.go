package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	zap *zap.Logger
	ctx context.Context
}

func New(ctx context.Context, cfg *Config) *Logger {
	loggerConfig := zap.NewProductionConfig()
	logLevel, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		panic(err)
	}
	loggerConfig.Level = logLevel

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	var ws zapcore.WriteSyncer
	if cfg.Async != nil {
		// Buffered, asynchronous writer
		ws = &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(os.Stdout),
			FlushInterval: cfg.Async.FlushInterval,
			Size:          cfg.Async.BufferSize,
		}
	} else {
		// Standard synchronous writer
		ws = zapcore.AddSync(os.Stdout)
	}

	// Create the core
	core := zapcore.NewCore(encoder, ws, logLevel)

	// Build the logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	defer logger.Sync()
	return &Logger{logger, ctx}
}

func (l *Logger) SetContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}
