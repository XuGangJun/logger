package mylogger

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

// 往文件里面写日志相关代码
// 对于Error类型的日志，可以考虑打印多一份到 errFileObj
type FileLogger struct {
	Level       LogLever
	filePath    string // 日志文件保存的路径
	fileName    string // 日志文件保存的文件名
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64 // 日志文件保存的最大内存
}

// NewFileLogger 构造器
func NewFileLogger(fp, fn string, maxSize int64) *FileLogger {
	f1 := &FileLogger{
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
	}
	err := f1.initFile()
	if err != nil {
		panic(err)
	}
	return f1
}

// 初始化日志文件
func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err1 := os.OpenFile(fullFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err1 != nil {
		fmt.Printf("open log file failed! err:%v", err1)
		return err1
	}
	errFileObj, err2 := os.OpenFile(fullFileName+"Error", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err2 != nil {
		fmt.Printf("open err log file failed! err:%v", err2)
		return err2
	}
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}

// 关闭文件操作
func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}

// 判断当前文件是否需要进行切割
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed! err:%v", err)
		return false
	}
	// 看下当前的日志文件是否 大于或等于 日志文件规定的max大小
	return fileInfo.Size() >= f.maxFileSize
}

// 对文件进行分割
func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	// 需要进行文件分割
	nowStr := time.Now().Format("20060102-150405.00")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed! err:%v\n", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name()) // 获取老的文件名
	newLogName := fmt.Sprintf("%s-%s", logName, nowStr)

	// 1. 关闭当前文件
	f.fileObj.Close()
	// 2. 备份当前文件
	os.Rename(logName, newLogName)
	// 3. 打开一个新的日志文件
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed! err:%v\n", err)
		return nil, err
	}
	// 4. 将打开的文件重新赋值给f.fileObj
	return fileObj, nil
}

// 记录日志的方法
func (f *FileLogger) log(lv LogLever, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		nowTime := time.Now().Format("2006-01-02 03:04:05")
		funcName, fileName, lineNo := getInfo(3)
		if f.checkSize(f.fileObj) { // 判断是否需要进行文件分割
			newFile, err := f.splitFile(f.fileObj) // 切割日志文件
			if err != nil {
				fmt.Printf("split file failed! err:%v\n", err)
				return
			}
			f.fileObj = newFile
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s\n", nowTime, getLogString(lv), fileName, funcName, lineNo, msg)
		if lv >= ERROR {
			// 源文件过大就要进行切割
			if f.checkSize(f.errFileObj) {
				newFile, err := f.splitFile(f.errFileObj) // 切割日志文件
				if err != nil {
					fmt.Printf("split file failed! err:%v\n", err)
					return
				}
				f.errFileObj = newFile
			}
			//如果要记录的日志大于ERROR了，再在错误日志里面记录一次
			fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %s\n", nowTime, getLogString(lv), fileName, funcName, lineNo, msg)
		}
	}
}

// 打印日志控制，默认全部打印
func (f *FileLogger) enable(loglever LogLever) bool {
	return loglever >= f.Level
}

func (f *FileLogger) Debug(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	f.log(DEBUG, format, a...)
}

func (f *FileLogger) Info(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	f.log(INFO, format, a...)
}

func (f *FileLogger) Warn(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	f.log(WARN, format, a...)
}

func (f *FileLogger) Error(format string, waitGroup *sync.WaitGroup, a ...interface{}) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}
	f.log(ERROR, format, a...)
}
