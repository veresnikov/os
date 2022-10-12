package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Info(message string)
	Error(err error) error
	Object(message string, object interface{})
}

func NewLogger(logrusLogger *logrus.Logger) Logger {
	return &logger{log: logrusLogger}
}

type logger struct {
	log *logrus.Logger
}

func (l *logger) Info(message string) {
	l.log.Info(message)
}

func (l *logger) Error(err error) error {
	if err == nil {
		return nil
	}
	l.log.Error(err)
	return err
}

func (l *logger) Object(message string, object interface{}) {
	l.log.Infof("message: %v; object: %v", message, object)
}
