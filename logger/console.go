package mylogger

import (
	"fmt"
	"sync"
	"time"
)

// Logger 日志结构体
type ConsoleLogger struct {
	Lever LogLever
}

func NewConsoleLog() ConsoleLogger {
	return ConsoleLogger{
		Lever: Undefined,
	}
}

func (c ConsoleLogger) log(lv LogLever, format string, a ...interface{}) {
	if c.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		nowTime := time.Now().Format("2006-01-02 03:04:05")
		funcName, fileName, lineNo := getInfo(3)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", nowTime, getLogString(lv), fileName, funcName, lineNo, msg)
	}
}

// 默认全部等级打印
func (c ConsoleLogger) enable(loglever LogLever) bool {
	return loglever >= c.Lever
}

func (c ConsoleLogger) Debug(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	c.log(DEBUG, format, a...)
}

func (c ConsoleLogger) Info(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	c.log(INFO, format, a...)
}

func (c ConsoleLogger) Warn(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	c.log(WARN, format, a...)
}

func (c ConsoleLogger) Error(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	c.log(ERROR, format, a...)
}
