package core

import (
	"fmt"
	"os"
	"runtime"
)

func Debug(args ...interface{}) {
	fmt.Println(append([]interface{}{"[DEBUG]"}, args...))
}

func Info(args ...interface{}) {
	fmt.Println(append([]interface{}{"[INFO]"}, args...))
}

func Warning(args ...interface{}) {
	fmt.Println(append([]interface{}{"[WARNING]"}, args...))
}

func Error(args ...interface{}) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	fmt.Println(append([]interface{}{"[ERROR]"}, args...))
	fmt.Println(string(buf))
	os.Exit(1)
}
