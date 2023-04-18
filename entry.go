package logrus

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

// Defines the key when adding errors using WithError.
var ErrorKey = "error"

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type Entry struct {
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	//filename:line
	//Location string
	FileName string

	Line int

	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, give a little extra room
		Data: make(Fields, 5),
	}
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(depth int, level Level, msg string) {
	var buffer *bytes.Buffer
	entry.Time = time.Now().UTC()
	entry.Level = level
	entry.Message = msg

	_, file, line, ok := runtime.Caller(2 + depth)
	if !ok {
		entry.FileName = "???"
		entry.Line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			entry.FileName = file[slash+1:]
		}
		entry.Line = line
	}
	//entry.Location = fmt.Sprintf("%s:%d", file, line)

	if err := entry.Logger.Hooks.Fire(level, &entry); err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
		entry.Logger.mu.Unlock()
	}
	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer
	serialized, err := entry.Logger.Formatter.Format(&entry)
	entry.Buffer = nil
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mu.Unlock()
	} else {
		entry.Logger.mu.Lock()
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
		entry.Logger.mu.Unlock()
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(0, DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(0, InfoLevel, fmt.Sprint(args...))
	}
}
func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(0, WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(0, ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(0, FatalLevel, fmt.Sprint(args...))
	}
	Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(0, PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

//Entry Ex family functions

func (entry *Entry) DebugEx(depth int, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(depth, DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) InfoEx(depth int, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(depth, InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) WarnEx(depth int, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(depth, WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) ErrorEx(depth int, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(depth, ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) FatalEx(depth int, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(depth, FatalLevel, fmt.Sprint(args...))
	}
	Exit(1)
}

func (entry *Entry) PanicEx(depth int, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(depth, PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.DebugEx(1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.InfoEx(1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.WarnEx(1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.ErrorEx(1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.FatalEx(1, fmt.Sprintf(format, args...))
	}
	Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.PanicEx(1, fmt.Sprintf(format, args...))
	}
}

//Entry PrintExf family functions
func (entry *Entry) DebugExf(depth int, format string, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.DebugEx(depth+1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) InfoExf(depth int, format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.InfoEx(depth+1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) WarnExf(depth int, format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.WarnEx(depth+1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) ErrorExf(depth int, format string, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.ErrorEx(depth+1, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) FatalExf(depth int, format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.FatalEx(depth+1, fmt.Sprintf(format, args...))
	}
	Exit(1)
}

func (entry *Entry) PanicExf(depth int, format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.PanicEx(1+depth, fmt.Sprintf(format, args...))
	}
}
