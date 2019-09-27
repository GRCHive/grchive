package core

import (
	"github.com/golang/glog"
	"os"
)

func Info(args ...interface{}) {
	glog.Infoln(args)
}

func Warning(args ...interface{}) {
	glog.Warningln(args)
}

func Error(args ...interface{}) {
	glog.Errorln(args)
	os.Exit(1)
}
