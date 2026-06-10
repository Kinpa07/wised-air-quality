package database

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	logger2 "github.com/SintroSecurity/go-libraries/logger"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

// Config logger config
type Config struct {
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  gormLogger.LogLevel
}

var (
	// Discard Discard logger will print any log to io.Discard
	Discard = New(context.TODO(), log.New(io.Discard, "", log.LstdFlags), Config{})
	// Default Default logger
	Default = New(context.Background(), log.New(os.Stdout, "\r\n", log.LstdFlags), Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gormLogger.Info,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      true,
		Colorful:                  true,
	})
)

// New initialize logger
func New(ctx context.Context, writer Writer, config Config) gormLogger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &logger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logger struct {
	Writer
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	jsonLogger := logger2.GetLoggerFromContext(ctx)
	jsonLogger.Info(msg, jsonLogger.Any("errors", data))
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	jsonLogger := logger2.GetLoggerFromContext(ctx)
	jsonLogger.Warn(msg, jsonLogger.Any("errors", data))
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	jsonLogger := logger2.GetLoggerFromContext(ctx)
	jsonLogger.Error(msg, jsonLogger.Any("errors", data))
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	jsonLogger := logger2.GetLoggerFromContext(ctx)

	sql, rows := fc()

	var notFound bool
	var success bool

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		notFound = true
	}

	duration := time.Since(begin)

	if err == nil {
		success = true
		jsonLogger.Debug(
			"sql tracing",
			jsonLogger.Any("begin_time", begin),
			jsonLogger.Any("time", time.Now()),
			jsonLogger.Int64("duration_ms", duration.Milliseconds()),
			jsonLogger.Any("sql", sql),
			jsonLogger.Any("rows", rows),
			jsonLogger.Any("success", success),
			jsonLogger.Any("not_found", notFound),
		)
	} else {
		jsonLogger.Error(
			"sql tracing",
			jsonLogger.Any("begin_time", begin),
			jsonLogger.Int64("duration_ms", duration.Milliseconds()),
			jsonLogger.Any("time", time.Now()),
			jsonLogger.Any("sql", sql),
			jsonLogger.Any("rows", rows),
			jsonLogger.Any("error", err),
			jsonLogger.Any("success", success),
		)
	}

	if metricsCollector.queryDuration != nil {
		//TODO doesn't generate real metrics with nested SQL queries. Should be fixed
		metricsCollector.queryDuration.WithLabelValues(getClause(sql), strconv.FormatBool(success)).Observe(float64(duration.Milliseconds()))
	}
}

// Trace print sql message
func (l logger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
