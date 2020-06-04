package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Format формат лога
type Format byte

const (
	// Text текстовый формат
	Text Format = iota
	// JSON формат JSON
	JSON
)

// Level уровень лога
type Level byte

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

// ILogger интерфейс логгера
type ILogger interface {
	// Print записывает в лог сообщение
	Print(level Level, args ...interface{})
	// Printf записывает в лог форматированное сообщение
	Printf(level Level, format string, args ...interface{})
}

// Logger логгер
type Logger struct {
	output io.Writer
	logger *logrus.Logger
}

// Option тип параметра логгера
type Option func(logger *Logger)

// New возвращает экземпляр логгера
func New(options ...Option) *Logger {
	logger := &Logger{
		output: nil,
		logger: logrus.New(),
	}

	for _, option := range options {
		option(logger)
	}

	return logger
}

// WithLevel устанавливает уровень логирования
func WithLevel(level Level) Option {
	return func(logger *Logger) {
		var l = logrus.InfoLevel

		if value, ok := levelMap[level]; ok {
			l = value
		}

		logger.logger.SetLevel(l)
	}
}

// WithFormat устанавливает формат лога
func WithFormat(format Format) Option {
	return func(logger *Logger) {
		switch format {
		case Text:
			logger.logger.SetFormatter(&logrus.TextFormatter{})
		case JSON:
			logger.logger.SetFormatter(&logrus.JSONFormatter{})
		default:
			logger.logger.SetFormatter(&logrus.TextFormatter{})
		}
	}
}

// WithOutput устанавливает вывод лога
func WithOutput(output io.Writer) Option {
	return func(logger *Logger) {
		logger.output = output

		if output != nil {
			logger.logger.SetOutput(output)
		}
	}
}

// Print записывает в лог сообщение
func (logger *Logger) Print(level Level, args ...interface{}) {
	if logger.output == nil {
		return
	}

	switch level {
	case TraceLevel:
		logger.logger.Trace(args...)
	case DebugLevel:
		logger.logger.Debug(args...)
	case InfoLevel:
		logger.logger.Info(args...)
	case WarningLevel:
		logger.logger.Warn(args...)
	case ErrorLevel:
		logger.logger.Error(args...)
	case FatalLevel:
		logger.logger.Fatal(args...)
	case PanicLevel:
		logger.logger.Panic(args...)
	}
}

// Printf записывает в лог форматированное сообщение
func (logger *Logger) Printf(level Level, format string, args ...interface{}) {
	if logger.output == nil {
		return
	}

	switch level {
	case TraceLevel:
		logger.logger.Tracef(format, args...)
	case DebugLevel:
		logger.logger.Debugf(format, args...)
	case InfoLevel:
		logger.logger.Infof(format, args...)
	case WarningLevel:
		logger.logger.Warnf(format, args...)
	case ErrorLevel:
		logger.logger.Errorf(format, args...)
	case FatalLevel:
		logger.logger.Fatalf(format, args...)
	case PanicLevel:
		logger.logger.Panicf(format, args...)
	}
}

var levelMap = map[Level]logrus.Level{
	TraceLevel:   logrus.TraceLevel,
	DebugLevel:   logrus.DebugLevel,
	InfoLevel:    logrus.InfoLevel,
	WarningLevel: logrus.WarnLevel,
	ErrorLevel:   logrus.ErrorLevel,
	FatalLevel:   logrus.FatalLevel,
	PanicLevel:   logrus.PanicLevel,
}
