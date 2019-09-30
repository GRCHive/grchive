package core

import (
	"k8s.io/klog"
	"os"
	"runtime"
)

func Info(args ...interface{}) {
	klog.Infoln(args)
}

func Warning(args ...interface{}) {
	klog.Warningln(args)
}

func Error(args ...interface{}) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	klog.Errorln(args)
	klog.Errorln(string(buf))
	os.Exit(1)
}
