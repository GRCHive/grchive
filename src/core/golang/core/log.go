package core

import (
	"k8s.io/klog"
	"os"
	"runtime"
)

func Debug(args ...interface{}) {
	klog.Infoln(append([]interface{}{"[DEBUG]"}, args...))
}

func Info(args ...interface{}) {
	klog.Infoln(append([]interface{}{"[INFO]"}, args...))
}

func Warning(args ...interface{}) {
	klog.Warningln(append([]interface{}{"[WARNING]"}, args...))
}

func Error(args ...interface{}) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	klog.Errorln(append([]interface{}{"[ERROR]"}, args...))
	klog.Errorln(string(buf))
	os.Exit(1)
}
