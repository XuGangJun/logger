package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
)

// LogLever 日志等级类型
type LogLever int

// 日志类型 等级变量
const (
	Undefined LogLever = iota
	DEBUG
	INFO
	WARN
	ERROR
)

// 将字符串转换成为日志等级对象   "INFO"  -->  INFO(LogLever)
func parseLogLever(s string) (LogLever, error) {
	s = strings.ToUpper(s)
	switch s {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		err := errors.New("未定义日志类型")
		return Undefined, err
	}
}

// 将日志等级对象转为字符串对象   INFO(LogLever)  -->  "INFO"
func getLogString(lv LogLever) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	default:
		return "DEBUG"
	}
}

// 获取文件信息
// runtime.Caller(skip)
// 调用者报告有关函数调用的文件和行号信息
// 调用 goroutine 的堆栈。参数 skip 是堆栈帧的数量
// 上升，0 标识 Caller 的调用者。 （由于历史原因
// 调用者和调用者之间跳过的含义不同。）返回值报告
// 对应的文件中的程序计数器、文件名和行号
// 称呼。如果无法恢复信息，则布尔值 ok 为 false。
func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller failed")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	funcName = strings.Split(funcName, ".")[1]
	lineNo = line
	return
}

// Logger 接口
type Logger interface {
	Debug(format string, waitGroup *sync.WaitGroup, a ...interface{})
	Info(format string, waitGroup *sync.WaitGroup, a ...interface{})
	Warn(format string, waitGroup *sync.WaitGroup, a ...interface{})
	Error(format string, waitGroup *sync.WaitGroup, a ...interface{})
}
