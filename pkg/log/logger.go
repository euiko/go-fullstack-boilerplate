package log

import (
	"context"
	"errors"
	"time"
)

type Level int

const (
	FatalLevel Level = iota
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type key int

var (
	loggerContextKey key
	fieldsContextKey key
	ErrNoLogger      = errors.New("no logger available")
)

type (
	// Logger represent any logging capable
	Logger interface {
		Log(level Level, msg *Log)
	}

	Fields map[string]interface{}

	// Log holds the data that will be logged
	Log struct {
		ctx     context.Context
		ts      time.Time
		message string
		fields  Fields
		err     error
	}

	Option interface {
		Configure(o *Log)
	}

	OptionFunc func(o *Log)
)

var globalLogger Logger

func (f OptionFunc) Configure(o *Log) {
	f(o)
}

// WithFileds adds all parameter fields for decorate current message option
func WithFields(fields Fields) Option {
	return OptionFunc(func(log *Log) {
		for k, v := range fields {
			log.fields[k] = v
		}
	})
}

// WithField add a single field to the message option
func WithField(key string, value interface{}) Option {
	return OptionFunc(func(log *Log) {
		log.fields[key] = value
	})
}

// WithContext decorate current message option to include a context
// maybe override the logger option
func WithContext(ctx context.Context) Option {
	return OptionFunc(func(log *Log) {
		log.ctx = ctx
	})
}

// WithTime decorate message option that override default now timestamp
func WithTime(t time.Time) Option {
	return OptionFunc(func(log *Log) {
		log.ts = t
	})
}

// WithError decorate message option that override the error
func WithError(err error) Option {
	return OptionFunc(func(log *Log) {
		log.err = err
	})
}

func Fatal(msg string, opts ...Option) error {
	return log(FatalLevel, msg, opts...)
}

func Error(msg string, opts ...Option) error {
	return log(ErrorLevel, msg, opts...)
}

func Warning(msg string, opts ...Option) error {
	return log(WarningLevel, msg, opts...)
}

func Info(msg string, opts ...Option) error {
	return log(InfoLevel, msg, opts...)
}

func Debug(msg string, opts ...Option) error {
	return log(DebugLevel, msg, opts...)
}

func Trace(msg string, opts ...Option) error {
	return log(TraceLevel, msg, opts...)
}

func SetDefault(logger Logger) {
	globalLogger = logger
}

func Default() Logger {
	return globalLogger
}

func SetFieldsContext(ctx context.Context, fields Fields) context.Context {
	// load or create fields from context
	current := getFieldsContext(ctx)
	if current == nil {
		current = make(Fields)
	}

	// merge all fields
	for k, v := range fields {
		current[k] = v
	}

	// inject to context
	return context.WithValue(ctx, loggerContextKey, current)
}

func log(level Level, msg string, opts ...Option) error {
	// global logger not yet specified, then do nothing
	if globalLogger == nil {
		return nil
	}

	log := newLog(msg, opts...)
	// TODO: support logger other than global
	logger := globalLogger

	// load all context specific option
	if log.ctx != nil {
		ctx := log.ctx

		// add additional fields defined in context
		// override all previously defined option
		fields := getFieldsContext(ctx)
		for k, v := range fields {
			log.fields[k] = v
		}
	}

	// do log
	logger.Log(level, log)
	return nil
}

func newLog(message string, options ...Option) *Log {

	// instantiate log
	log := &Log{
		ts:      time.Now(),
		message: message,
		fields:  make(Fields),
		ctx:     nil,
		err:     nil,
	}

	// load all decorator function
	for _, o := range options {
		o.Configure(log)
	}

	return log
}

func getFieldsContext(ctx context.Context) Fields {
	instance := ctx.Value(fieldsContextKey)
	if instance == nil {
		return nil
	}

	fields, ok := instance.(Fields)
	if !ok {
		return nil
	}

	return fields
}

func init() {
	// set the default logger to use logrus
	SetDefault(NewLogrusLogger(InfoLevel))
}
