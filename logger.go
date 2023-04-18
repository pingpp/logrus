package logrus

import (
	"io"
	"os"
	"sync"
)

type Logger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventorous, such as logging to Kafka.
	Out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks LevelHooks
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter
	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged. `logrus.Debug` is useful in
	Level Level
	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

// Creates a new logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance. You can also just
// instantiate your own:
//
//    var log = &Logger{
//      Out: os.Stderr,
//      Formatter: new(JSONFormatter),
//      Hooks: make(LevelHooks),
//      Level: logrus.DebugLevel,
//    }
//
// It's recommended to make this a global instance called `log`.
func New() *Logger {
	return &Logger{
		Out:       os.Stdout,
		Formatter: new(TextFormatter),
		Hooks:     make(LevelHooks),
		Level:     InfoLevel,
	}
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	logger.entryPool.Put(entry)
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

//logger Print family
func (logger *Logger) Debug(args ...interface{}) {
	if logger.Level >= DebugLevel {
		entry := logger.newEntry()
		entry.DebugEx(1, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Info(args ...interface{}) {
	if logger.Level >= InfoLevel {
		entry := logger.newEntry()
		entry.InfoEx(1, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Warn(args ...interface{}) {
	if logger.Level >= WarnLevel {
		entry := logger.newEntry()
		entry.WarnEx(1, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if logger.Level >= ErrorLevel {
		entry := logger.newEntry()
		entry.ErrorEx(1, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.Level >= FatalLevel {
		entry := logger.newEntry()
		entry.FatalEx(1, args...)
		logger.releaseEntry(entry)
	}
	Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	if logger.Level >= PanicLevel {
		entry := logger.newEntry()
		entry.PanicEx(1, args...)
		logger.releaseEntry(entry)
	}
}

//logger PrintEx family
func (logger *Logger) DebugEx(depth int, args ...interface{}) {
	if logger.Level >= DebugLevel {
		entry := logger.newEntry()
		entry.DebugEx(1+depth, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) InfoEx(depth int, args ...interface{}) {
	if logger.Level >= InfoLevel {
		entry := logger.newEntry()
		entry.InfoEx(1+depth, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) WarnEx(depth int, args ...interface{}) {
	if logger.Level >= WarnLevel {
		entry := logger.newEntry()
		entry.WarnEx(1+depth, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) ErrorEx(depth int, args ...interface{}) {
	if logger.Level >= ErrorLevel {
		entry := logger.newEntry()
		entry.ErrorEx(1+depth, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) FatalEx(depth int, args ...interface{}) {
	if logger.Level >= FatalLevel {
		entry := logger.newEntry()
		entry.FatalEx(1+depth, args...)
		logger.releaseEntry(entry)
	}
	Exit(1)
}

func (logger *Logger) PanicEx(depth int, args ...interface{}) {
	if logger.Level >= PanicLevel {
		entry := logger.newEntry()
		entry.PanicEx(1+depth, args...)
		logger.releaseEntry(entry)
	}
}

// logger Printf family functions
func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.Level >= DebugLevel {
		entry := logger.newEntry()
		entry.DebugExf(1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.Level >= InfoLevel {
		entry := logger.newEntry()
		entry.InfoExf(1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		entry := logger.newEntry()
		entry.WarnExf(1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		entry := logger.newEntry()
		entry.WarnExf(1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.Level >= ErrorLevel {
		entry := logger.newEntry()
		entry.ErrorExf(1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.Level >= FatalLevel {
		entry := logger.newEntry()
		entry.FatalExf(1, format, args...)
		logger.releaseEntry(entry)
	}
	Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.Level >= PanicLevel {
		entry := logger.newEntry()
		entry.PanicExf(1, format, args...)
		logger.releaseEntry(entry)
	}
}

//logger PrintExf family

func (logger *Logger) DebugExf(depth int, format string, args ...interface{}) {
	if logger.Level >= DebugLevel {
		entry := logger.newEntry()
		entry.DebugExf(depth+1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) InfoExf(depth int, format string, args ...interface{}) {
	if logger.Level >= InfoLevel {
		entry := logger.newEntry()
		entry.InfoExf(depth+1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) WarnExf(depth int, format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		entry := logger.newEntry()
		entry.WarnExf(depth+1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) WarningExf(depth int, format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		entry := logger.newEntry()
		entry.WarnExf(depth+1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) ErrorExf(depth int, format string, args ...interface{}) {
	if logger.Level >= ErrorLevel {
		entry := logger.newEntry()
		entry.ErrorExf(depth+1, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) FatalExf(depth int, format string, args ...interface{}) {
	if logger.Level >= FatalLevel {
		entry := logger.newEntry()
		entry.FatalExf(depth+1, format, args...)
		logger.releaseEntry(entry)
	}
	Exit(1)
}

func (logger *Logger) PanicExf(depth int, format string, args ...interface{}) {
	if logger.Level >= PanicLevel {
		entry := logger.newEntry()
		entry.PanicExf(depth+1, format, args...)
		logger.releaseEntry(entry)
	}
}

//When file is opened with appending mode, it's safe to
//write concurrently to a file (within 4k message on Linux).
//In these cases user can choose to disable the lock.
func (logger *Logger) SetNoLock() {
	logger.mu.Disable()
}
