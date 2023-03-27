package logging

type Logger interface {
	Print(level Level, message string) error
	Printf(level Level, format string, args ...any) error

	Error(message string) error
	Errorf(format string, args ...any) error

	Warn(message string) error
	Warnf(format string, args ...any) error

	Info(message string) error
	Infof(format string, args ...any) error

	Trace(message string) error
	Tracef(format string, args ...any) error

	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}
