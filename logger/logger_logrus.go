package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type myLogger struct {
	l   *logrus.Entry
	lev Level
}

func (m *myLogger) WithField(name string, value interface{}) Logger {
	return newLogger(m.l.WithField(name, value), m.lev)
}

func (m *myLogger) DebugF(format string, args ...interface{}) {
	if m.lev <= Debug {
		m.l.Debugf(format, args...)
	}
}

func (m *myLogger) InfoF(format string, args ...interface{}) {
	if m.lev <= Info {
		m.l.Infof(format, args...)
	}
}

func (m *myLogger) WarnF(format string, args ...interface{}) {
	if m.lev <= Warn {
		m.l.Warnf(format, args...)
	}
}

func (m *myLogger) ErrorF(format string, args ...interface{}) {
	if m.lev <= Error {
		m.l.Errorf(format, args...)
	}
}

func (m *myLogger) FatalF(format string, args ...interface{}) {
	if m.lev <= Fatal {
		m.l.Fatalf(format, args...)
	}
}

func (m *myLogger) SetLevel(level Level) {
	m.lev = level
}

func newLogger(entry *logrus.Entry, level Level) *myLogger {
	return &myLogger{l: entry, lev: level}
}

// Options used to config the log
type Options struct {
	Writer     io.Writer
	Level      Level
	JsonFormat bool
}

// Option used to change Options by closure
type Option func(options *Options)

// WithWriter is set the log writer
func WithWriter(writer io.Writer) Option {
	return func(options *Options) {
		options.Writer = writer
	}
}

// WithLevel is set the log level
func WithLevel(level Level) Option {
	return func(options *Options) {
		options.Level = level
	}
}

// WithJsonFormat is config whether to use json format log
func WithJsonFormat(useJsonFormat bool) Option {
	return func(options *Options) {
		options.JsonFormat = useJsonFormat
	}
}

// NewLogger is create a logger based by logrus implements Logger
func NewLogger(opts ...Option) *myLogger {
	options := Options{Level: Info, Writer: os.Stdout, JsonFormat: true}
	for _, o := range opts {
		o(&options)
	}
	l := logrus.StandardLogger()
	l.SetOutput(options.Writer)
	if options.JsonFormat {
		l.SetFormatter(new(logrus.JSONFormatter))
	} else {
		l.SetFormatter(new(logrus.TextFormatter))
	}
	return newLogger(logrus.NewEntry(l), options.Level)
}
