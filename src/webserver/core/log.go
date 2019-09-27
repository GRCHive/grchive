package core

import (
	"github.com/golang/glog"
	"os"
)

func Warning(args ...interface{}) {
	glog.Warningln(args)
}

func Error(args ...interface{}) {
	glog.Errorln(args)
	os.Exit(1)
}
