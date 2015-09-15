/*
Package glog wraps github.com/golang/glog for grpc-gateway logging.
*/
package glog

import (
	"github.com/gengo/grpc-gateway/log"
	"github.com/golang/glog"
)

func init() {
	log.SetLogger(&gLogger{})
}

type gLogger struct{}

func (l *gLogger) Infof(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func (l *gLogger) Infoln(args ...interface{}) {
	glog.Infoln(args...)
}

func (l *gLogger) Warnf(format string, args ...interface{}) {
	glog.Warningf(format, args...)
}

func (l *gLogger) Warnln(args ...interface{}) {
	glog.Warningln(args...)
}

func (l *gLogger) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}

func (l *gLogger) Errorln(args ...interface{}) {
	glog.Errorln(args...)
}
