package logger

// Logger defined Logger for log anything
type Logger interface {
	// WithField will config a field with name and value
	// This method will return a new Logger
	WithField(name string, value interface{}) Logger

	// DebugF log anything at debug level
	DebugF(format string, args ...interface{})

	// InfoF log anything at info level
	InfoF(format string, args ...interface{})

	// WarnF log anything at warning level
	WarnF(format string, args ...interface{})

	// ErrorF log anything at error level
	ErrorF(format string, args ...interface{})

	// FatalF log anything at fatal level
	FatalF(format string, args ...interface{})

	// SetLevel is set the Logger log level
	SetLevel(level Level)
}

// Level is the Log Level
type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
	Fatal Level = "fatal"
)
