package commonlogger

import (
	"github.com/johnmackenzie91/commonlogger/resolvers"

	"github.com/sirupsen/logrus"
)

// ErrorInfoDebugger represents the entire interface exposed by this package
type ErrorInfoDebugger interface {
	Error(...interface{})
	Info(...interface{})
	Debug(...interface{})
}

// Logger the struct which implements the ErrorInfoDebugger interface
type Logger struct {
	logger     *logrus.Logger
	strategies resolvers.Strategies
}

// New instansiates a new Logger
func New(logger *logrus.Logger, c Config) Logger {
	l := Logger{
		logger: logger,
	}
	if c.ContextResolver != nil {
		l.strategies.ContextResolver = c.ContextResolver
	}
	if c.RequestResolver != nil {
		l.strategies.RequestResolver = c.RequestResolver
	}
	if c.ResponseResolver != nil {
		l.strategies.ResponseResolver = c.ResponseResolver
	}
	return l
}

// Debug writes an error line
func (l Logger) Debug(opts ...interface{}) {
	line := l.buildLogLine(opts...)
	l.logger.WithTime(line.time).WithContext(line.ctx).WithFields(line.fields).Debug(line.msg)
}

// Info writes an info line
func (l Logger) Info(opts ...interface{}) {
	line := l.buildLogLine(opts...)
	l.logger.WithTime(line.time).WithContext(line.ctx).WithFields(line.fields).Info(line.msg)
}

// Error writes an error line
func (l Logger) Error(opts ...interface{}) {
	line := l.buildLogLine(opts...)
	l.logger.WithTime(line.time).WithContext(line.ctx).WithFields(line.fields).Error(line.msg)
}
