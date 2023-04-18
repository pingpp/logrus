package logrus

import (
	"io"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func StandardLogger() *Logger {
	return std
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Out = out
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Formatter = formatter
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Level = level
}

// GetLevel returns the standard logger level.
func GetLevel() Level {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.Level
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Hooks.Add(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	std.DebugEx(1, args...)
}

// // Print logs a message at level Info on the standard logger.
// func Print(args ...interface{}) {
// 	std.Print(args...)
// }

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	std.InfoEx(1, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	std.WarnEx(1, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	std.ErrorEx(1, args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	std.PanicEx(1, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	std.FatalEx(1, args...)
}

//PrintEx Family
// Debug logs a message at level Debug on the standard logger.
func DebugEx(depth int, args ...interface{}) {
	std.DebugEx(depth+1, args...)
}

// // Print logs a message at level Info on the standard logger.
// func Print(args ...interface{}) {
// 	std.Print(args...)
// }

// Info logs a message at level Info on the standard logger.
func InfoEx(depth int, args ...interface{}) {
	std.InfoEx(depth+1, args...)
}

// Warn logs a message at level Warn on the standard logger.
func WarnEx(depth int, args ...interface{}) {
	std.WarnEx(depth+1, args...)
}

// Error logs a message at level Error on the standard logger.
func ErrorEx(depth int, args ...interface{}) {
	std.ErrorEx(depth+1, args...)
}

// Panic logs a message at level Panic on the standard logger.
func PanicEx(depth int, args ...interface{}) {
	std.PanicEx(depth+1, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func FatalEx(depth int, args ...interface{}) {
	std.FatalEx(depth+1, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	std.DebugExf(1, format, args...)
}

// Printf logs a message at level Info on the standard logger.
// func Printf(format string, args ...interface{}) {
// 	std.Printf(format, args...)
// }

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	std.InfoExf(1, format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	std.WarnExf(1, format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
// func Warningf(format string, args ...interface{}) {
// 	std.WarningExf(1, format, args...)
// }

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	std.ErrorExf(1, format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	std.PanicExf(1, format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	std.FatalExf(1, format, args...)
}

func DebugExf(depth int, format string, args ...interface{}) {
	std.DebugExf(1+depth, format, args...)
}

// Printf logs a message at level Info on the standard logger.
// func Printf(format string, args ...interface{}) {
// 	std.Printf(format, args...)
// }

// Infof logs a message at level Info on the standard logger.
func InfoExf(depth int, format string, args ...interface{}) {
	std.InfoExf(1+depth, format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func WarnExf(depth int, format string, args ...interface{}) {
	std.WarnExf(1+depth, format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
// func Warningf(format string, args ...interface{}) {
// 	std.WarningExf(1, format, args...)
// }

// Errorf logs a message at level Error on the standard logger.
func ErrorExf(depth int, format string, args ...interface{}) {
	std.ErrorExf(1+depth, format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func PanicExf(depth int, format string, args ...interface{}) {
	std.PanicExf(1+depth, format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func FatalExf(depth int, format string, args ...interface{}) {
	std.FatalExf(1+depth, format, args...)
}
