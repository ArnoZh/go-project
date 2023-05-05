// Package util .
 
package util

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/sirupsen/logrus"
)

// PrintPanicStack recover并打印堆栈 用法: defer utils.PrintPanicStack();注意 defer func() {utils.PrintPanicStack()} 调用无效
func PrintPanicStack() {
	if r := recover(); r != nil {
		logrus.Errorf("%v: %s", r, debug.Stack())
	}
}

// PrintDebugStack 打印调用堆栈
func PrintDebugStack() {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		logrus.Debugf("%s:%d %s\n", file, line, f.Name())
	}
}

// Stack 返回调用堆栈
func Stack() string {
	s := strings.Builder{}
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		s.WriteString(fmt.Sprintf("%s:%d %s ", file, line, f.Name()))
	}
	return s.String()
}
