package main

import (
	mylogger "mylogger/logger"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	count := 2
	var consoleLog mylogger.Logger = mylogger.NewConsoleLog()
	var fileLogger mylogger.Logger = mylogger.NewFileLogger("main", "firstLog", 1024*1024)
	for i := 0; i < count; i++ {
		wg.Add(8)
		go consoleLog.Info("sample Info---a", &wg)
		go consoleLog.Debug("Debug Log test---b", &wg)
		go consoleLog.Warn("Warn Something---c", &wg)
		go consoleLog.Error("Something Error---d", &wg)

		go fileLogger.Info("file sample Info---1", &wg)
		go fileLogger.Debug("file Debug Log test---2", &wg)
		go fileLogger.Warn("file Warn Something---3", &wg)
		go fileLogger.Error("file Something Error---4", &wg)
	}
	wg.Wait()
}
