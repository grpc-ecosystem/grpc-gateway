package log

import "log"

type logger struct {
	*log.Logger
}

func newLogger(l *log.Logger) *logger {
	return &logger{l}
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.Println(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func (l *logger) Warnln(args ...interface{}) {
	l.Println(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.Println(args...)
}
