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
func (self *Logger) Print(level Level, args ...interface{}) {
	if self.output == nil {
		return
	}

	switch level {
	case TraceLevel:
		self.logger.Trace(args...)
	case DebugLevel:
		self.logger.Debug(args...)
	case InfoLevel:
		self.logger.Info(args...)
	case WarningLevel:
		self.logger.Warn(args...)
	case ErrorLevel:
		self.logger.Error(args...)
	case FatalLevel:
		self.logger.Fatal(args...)
	case PanicLevel:
		self.logger.Panic(args...)
	}
}

// Printf записывает в лог форматированное сообщение
func (self *Logger) Printf(level Level, format string, args ...interface{}) {
	if self.output == nil {
		return
	}

	switch level {
	case TraceLevel:
		self.logger.Tracef(format, args...)
	case DebugLevel:
		self.logger.Debugf(format, args...)
	case InfoLevel:
		self.logger.Infof(format, args...)
	case WarningLevel:
		self.logger.Warnf(format, args...)
	case ErrorLevel:
		self.logger.Errorf(format, args...)
	case FatalLevel:
		self.logger.Fatalf(format, args...)
	case PanicLevel:
		self.logger.Panicf(format, args...)
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
