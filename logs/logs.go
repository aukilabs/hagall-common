package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aukilabs/hagall-common/errors"
	"go.opentelemetry.io/otel/trace"
)

var (
	// The function to encode log entries and their tags.
	Encoder func(any) ([]byte, error)
)

// Set the function that logs an entry.
func SetLogger(l func(Entry)) {
	logger = l
	SetLevel(currentLevel)
}

// Available log levels.
const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrorLevel

	AppKeyTag        = "app_key"
	SessionIDTag     = "session_id"
	ParticipantIDTag = "participant_id"
	ClientIDTag      = "client_id"
)

// A log level.
type Level int

// Parses a level from a string.
func ParseLevel(v string) Level {
	switch v {
	case "debug":
		return DebugLevel

	case "info":
		return InfoLevel

	case "warning":
		return WarningLevel

	case "error":
		return ErrorLevel

	default:
		l, _ := strconv.Atoi(v)
		return Level(l)
	}
}

// Converts a level to a string.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"

	case InfoLevel:
		return "info"

	case WarningLevel:
		return "warning"

	case ErrorLevel:
		return "error"

	default:
		return strconv.Itoa(int(l))
	}
}

// Sets what log levels are logged. Levels under the given level are ignored.
func SetLevel(v Level) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	currentLevel = v

	for i := DebugLevel; i <= ErrorLevel; i++ {
		if i < v {
			loggers[i] = empytLogger
		} else {
			loggers[i] = logger
		}
	}
}

// SetInlineEncoder is a helper function that set the error encoder to
// json.Marshal.
func SetInlineEncoder() {
	Encoder = json.Marshal
}

// SetIndentEncoder is a helper function that set the error encoder to a
// function that uses json.MarshalIndent.
func SetIndentEncoder() {
	Encoder = func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}
}

// Creates a log entry.
func New() Entry {
	return entry{}
}

// Creates a log entry and set the given tag.
func WithTag(k string, v any) Entry {
	return New().WithTag(k, v)
}

// Creates a log entry with client id
func WithClientID(v string) Entry {
	return New().WithTag(ClientIDTag, v)
}

// Creates a log entry with opentelemetry context tags
func WithOtelCtx(ctx context.Context) Entry {
	return New().WithOtelCtx(ctx)
}

// Logs an error.
func Error(err error) {
	New().Error(err)
}

// Logs with warn severity
func Warn(arg ...any) {
	New().Warn(arg)
}

// Logs with warn severity
func Warnf(format string, arg ...any) {
	New().Warnf(format, arg)
}

// Logs with info severity
func Info(arg ...any) {
	New().Info(arg)
}

// Logs with info severity
func Infof(format string, arg ...any) {
	New().Infof(format, arg)
}

// Logs with debug severity
func Debug(arg ...any) {
	New().Debug(arg)
}

// Logs with debug severity
func Debugf(format string, arg ...any) {
	New().Debugf(format, arg)
}

// Logs with error severity and panic
func Panic(err error) {
	New().Panic(err)
}

// Logs with error severity and exit with status code 1
func Fatal(err error) {
	New().Fatal(err)
}

type Entry interface {
	// Return the time when the entry was created.
	Time() time.Time

	// Returns the log level.
	Level() Level

	// Sets the tag key with the given value. The value is converted to a
	// string.
	WithTag(k string, v any) Entry

	// Set Client ID tag
	WithClientID(v string) Entry

	// Set opentelemetry context tags
	WithOtelCtx(ctx context.Context) Entry

	// Returns the log tags.
	Tags() map[string]any

	// Logs the given values with debug level.
	Debug(v ...any)

	// Logs the velues with the given format on with debug level.
	Debugf(format string, v ...any)

	// Logs the given values with info level.
	Info(v ...any)

	// Logs the velues with the given format on with info level.
	Infof(format string, v ...any)

	// Logs the given values with warning level.
	Warn(v ...any)

	// Logs the velues with the given format on with warning level.
	Warnf(format string, v ...any)

	// Logs the given values with error level.
	Error(error)

	// Logs error on error level and exit with status 1.
	Fatal(err error)

	// Logs error on error level and panic
	Panic(err error)

	// Returns the error used to create the entry.
	GetError() error

	// Return the entry as a string.
	String() string
}

var (
	loggerMutex   sync.RWMutex
	loggers       = make(map[Level]func(Entry), ErrorLevel+1)
	logger        func(e Entry)
	defaultLogger = func(e Entry) { fmt.Println(e) }
	empytLogger   = func(Entry) {}
	currentLevel  Level
)

func init() {
	SetInlineEncoder()
	SetLogger(defaultLogger)
}

func log(e Entry) {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	loggers[e.Level()](e)
}

type entry struct {
	time    time.Time
	level   Level
	message string
	tags    map[string]any
	err     error
}

func (e entry) Time() time.Time {
	return e.time
}

func (e entry) Level() Level {
	return e.level
}

func (e entry) WithTag(k string, v any) Entry {
	if e.tags == nil {
		e.tags = make(map[string]any)
	}

	e.tags[k] = normalizeTag(v)
	return e
}

func (e entry) WithClientID(v string) Entry {
	return e.WithTag(ClientIDTag, v)
}

func (e entry) WithOtelCtx(ctx context.Context) Entry {
	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID()
	spanID := trace.SpanFromContext(ctx).SpanContext().SpanID()

	if traceID.IsValid() && spanID.IsValid() {
		return e.WithTag("trace-id", traceID.String()).WithTag("span-id", spanID.String())
	}
	return e
}

func (e entry) Tags() map[string]any {
	return e.tags
}

func (e entry) Debug(v ...any) {
	e.time = time.Now()
	e.level = DebugLevel
	e.message = fmt.Sprint(v...)
	log(e)
}

func (e entry) Debugf(format string, v ...any) {
	e.time = time.Now()
	e.level = DebugLevel
	e.message = fmt.Sprintf(format, v...)
	log(e)
}

func (e entry) Info(v ...any) {
	e.time = time.Now()
	e.level = InfoLevel
	e.message = fmt.Sprint(v...)
	log(e)
}

func (e entry) Infof(format string, v ...any) {
	e.time = time.Now()
	e.level = InfoLevel
	e.message = fmt.Sprintf(format, v...)
	log(e)
}

func (e entry) Warn(v ...any) {
	e.time = time.Now()
	e.level = WarningLevel
	e.message = fmt.Sprint(v...)
	log(e)
}

func (e entry) Warnf(format string, v ...any) {
	e.time = time.Now()
	e.level = WarningLevel
	e.message = fmt.Sprintf(format, v...)
	log(e)
}

func (e entry) Error(err error) {
	e.time = time.Now()
	e.level = ErrorLevel
	e.message = errors.Message(err)
	e.err = err
	log(e)
}

func (e entry) Fatal(err error) {
	e.time = time.Now()
	e.level = ErrorLevel
	e.message = errors.Message(err)
	e.err = err
	log(e)
	os.Exit(1)
}

func (e entry) Panic(err error) {
	e.time = time.Now()
	e.level = ErrorLevel
	e.message = errors.Message(err)
	e.err = err
	log(e)
	panic(e)
}

func (e entry) GetError() error {
	return e.err
}

func (e entry) String() string {
	b, _ := e.MarshalJSON()
	return string(b)
}

func (e entry) MarshalJSON() ([]byte, error) {
	var line string
	var typ string
	var wrappedErr error

	if err, ok := e.err.(errors.Error); ok {
		line = err.Line()
		typ = err.Type()
		wrappedErr = err.Unwrap()

		if e.tags == nil && len(err.Tags()) != 0 {
			e.tags = make(map[string]any)
		}
		for k, v := range err.Tags() {
			e.tags[k] = v
		}
	}

	return Encoder(struct {
		Time    time.Time      `json:"time"`
		Level   string         `json:"level"`
		Message string         `json:"message"`
		Line    string         `json:"line,omitempty"`
		Type    string         `json:"type,omitempty"`
		Tags    map[string]any `json:"tags,omitempty"`
		Wrap    error          `json:"wrap,omitempty"`
	}{
		Time:    e.time,
		Level:   e.level.String(),
		Message: e.message,
		Tags:    e.tags,
		Line:    line,
		Type:    typ,
		Wrap:    wrappedErr,
	})
}

func normalizeTag(v any) any {
	switch v := v.(type) {
	case error:
		return v.Error()

	case time.Duration:
		return v.String()

	case []byte:
		return string(v)

	default:
		return v
	}
}
