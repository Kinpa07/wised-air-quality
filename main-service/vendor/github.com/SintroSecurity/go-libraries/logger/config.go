package logger

import "time"

type Config struct {
	Level string `mapstructure:"LEVEL"`
	Async *Async `mapstructure:"ASYNC"`
}

type Async struct {
	FlushInterval time.Duration `mapstructure:"FLUSHINTERVAL"`
	BufferSize    int           `mapstructure:"BUFFERSIZE"`
}
