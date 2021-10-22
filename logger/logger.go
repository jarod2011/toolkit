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
type Level int

const (
	Debug Level = iota - 1
	Info
	Warn
	Error
	Fatal
)

// String is description the log leve name
func (l Level) String() string {
	return []string{"debug", "info", "warning", "error", "fatal"}[int(l)+1]
}
