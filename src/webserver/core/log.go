package core

import (
	"k8s.io/klog"
	"os"
)

func Info(args ...interface{}) {
	klog.Infoln(args)
}

func Warning(args ...interface{}) {
	klog.Warningln(args)
}

func Error(args ...interface{}) {
	klog.Errorln(args)
	os.Exit(1)
}
