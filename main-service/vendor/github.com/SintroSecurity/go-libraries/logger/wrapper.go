package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"unsafe"
)

type Field zap.Field

const reqIDFieldKey = "request_id"

func (l *Logger) With(fields ...Field) *Logger {
	// Convert your wrapper Field to zap.Field
	zapFields := unsafe.Slice((*zap.Field)(unsafe.Pointer(&fields[0])), len(fields))

	return &Logger{
		zap: l.zap.With(zapFields...),
		ctx: l.ctx,
	}
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, *(*[]zap.Field)(unsafe.Pointer(&fields))...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, *(*[]zap.Field)(unsafe.Pointer(&fields))...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, *(*[]zap.Field)(unsafe.Pointer(&fields))...)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, *(*[]zap.Field)(unsafe.Pointer(&fields))...)
}

func (l *Logger) Binary(key string, val []byte) Field {
	return Field(zap.Binary(key, val))
}

func (l *Logger) Bool(key string, val bool) Field {
	return Field(zap.Bool(key, val))
}

func (l *Logger) Boolp(key string, val *bool) Field {
	return Field(zap.Boolp(key, val))
}

func (l *Logger) ByteString(key string, val []byte) Field {
	return Field(zap.ByteString(key, val))
}

func (l *Logger) Complex128(key string, val complex128) Field {
	return Field(zap.Complex128(key, val))
}

func (l *Logger) Complex128p(key string, val *complex128) Field {
	return Field(zap.Complex128p(key, val))
}

func (l *Logger) Complex64(key string, val complex64) Field {
	return Field(zap.Complex64(key, val))
}

func (l *Logger) Complex64p(key string, val *complex64) Field {
	return Field(zap.Complex64p(key, val))
}

func (l *Logger) Err(err error) Field {
	return Field(zap.Error(err))
}

func (l *Logger) Float64(key string, val float64) Field {
	return Field(zap.Float64(key, val))
}

func (l *Logger) Float64p(key string, val *float64) Field {
	return Field(zap.Float64p(key, val))
}

func (l *Logger) Float32(key string, val float32) Field {
	return Field(zap.Float32(key, val))
}

func (l *Logger) Float32p(key string, val *float32) Field {
	return Field(zap.Float32p(key, val))
}

func (l *Logger) Int(key string, val int) Field {
	return Field(zap.Int(key, val))
}

func (l *Logger) Intp(key string, val *int) Field {
	return Field(zap.Intp(key, val))
}

func (l *Logger) Int64(key string, val int64) Field {
	return Field(zap.Int64(key, val))
}

func (l *Logger) Int64p(key string, val *int64) Field {
	return Field(zap.Int64p(key, val))
}

func (l *Logger) Int32(key string, val int32) Field {
	return Field(zap.Int32(key, val))
}

func (l *Logger) Int32p(key string, val *int32) Field {
	return Field(zap.Int32p(key, val))
}

func (l *Logger) Int16(key string, val int16) Field {
	return Field(zap.Int16(key, val))
}

func (l *Logger) Int16p(key string, val *int16) Field {
	return Field(zap.Int16p(key, val))
}

func (l *Logger) Int8(key string, val int8) Field {
	return Field(zap.Int8(key, val))
}

func (l *Logger) Int8p(key string, val *int8) Field {
	return Field(zap.Int8p(key, val))
}

func (l *Logger) String(key string, val string) Field {
	return Field(zap.String(key, val))
}

func (l *Logger) Stringp(key string, val *string) Field {
	return Field(zap.Stringp(key, val))
}

func (l *Logger) Uint(key string, val uint) Field {
	return Field(zap.Uint(key, val))
}

func (l *Logger) Uintp(key string, val *uint) Field {
	return Field(zap.Uintp(key, val))
}

func (l *Logger) Uint64(key string, val uint64) Field {
	return Field(zap.Uint64(key, val))
}

func (l *Logger) Uint64p(key string, val *uint64) Field {
	return Field(zap.Uint64p(key, val))
}

func (l *Logger) Uint32(key string, val uint32) Field {
	return Field(zap.Uint32(key, val))
}

func (l *Logger) Uint32p(key string, val *uint32) Field {
	return Field(zap.Uint32p(key, val))
}

func (l *Logger) Uint16(key string, val uint16) Field {
	return Field(zap.Uint16(key, val))
}

func (l *Logger) Uint16p(key string, val *uint16) Field {
	return Field(zap.Uint16p(key, val))
}

func (l *Logger) Uint8(key string, val uint8) Field {
	return Field(zap.Uint8(key, val))
}

func (l *Logger) Uint8p(key string, val *uint8) Field {
	return Field(zap.Uint8p(key, val))
}

func (l *Logger) Uintptr(key string, val uintptr) Field {
	return Field(zap.Uintptr(key, val))
}

func (l *Logger) Uintptrp(key string, val *uintptr) Field {
	return Field(zap.Uintptrp(key, val))
}

func (l *Logger) Reflect(key string, val interface{}) Field {
	return Field(zap.Reflect(key, val))
}

func (l *Logger) Namespace(key string) Field {
	return Field(zap.Namespace(key))
}

func (l *Logger) Stringer(key string, val fmt.Stringer) Field {
	return Field(zap.Stringer(key, val))
}

func (l *Logger) Time(key string, val time.Time) Field {
	return Field(zap.Time(key, val))
}

func (l *Logger) Timep(key string, val *time.Time) Field {
	return Field(zap.Timep(key, val))
}

func (l *Logger) Stack(key string) Field {
	return Field(zap.Stack(key))
}

func (l *Logger) StackSkip(key string, skip int) Field {
	return Field(zap.StackSkip(key, skip))
}

func (l *Logger) Duration(key string, val time.Duration) Field {
	return Field(zap.Duration(key, val))
}

func (l *Logger) Durationp(key string, val *time.Duration) Field {
	return Field(zap.Durationp(key, val))
}

func (l *Logger) Object(key string, val zapcore.ObjectMarshaler) Field {
	return Field(zap.Object(key, val))
}

func (l *Logger) Inline(val zapcore.ObjectMarshaler) Field {
	return Field(zap.Inline(val))
}

func (l *Logger) Any(key string, value interface{}) Field {
	return Field(zap.Any(key, value))
}
