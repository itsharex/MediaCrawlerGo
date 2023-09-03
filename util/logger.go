package util

import (
	"fmt"
	"os"
	"time"
)

const (
	// LevelError 错误
	LevelError = iota
	// LevelWarning 警告
	LevelWarning
	// LevelInformational 提示
	LevelInformational
	// LevelDebug 除错
	LevelDebug
)

var logger *Logger

// Logger 日志
type Logger struct {
	name  string
	level int
}

// Println 打印
func (ll *Logger) Println(msg string) {
	fmt.Printf("%s %s %s \n", time.Now().Format("2006-01-02 15:04:05 -0700"), ll.name, msg)
}

// Panic 极端错误
func (ll *Logger) Panic(format string, v ...any) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf(ll.name+"[Panic] "+format, v...)
	ll.Println(msg)
	os.Exit(0)
}

// Error 错误
func (ll *Logger) Error(format string, v ...any) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf("[E] "+format, v...)
	ll.Println(msg)
}

// Warning 警告
func (ll *Logger) Warning(format string, v ...any) {
	if LevelWarning > ll.level {
		return
	}
	msg := fmt.Sprintf("[W] "+format, v...)
	ll.Println(msg)
}

// Info 信息
func (ll *Logger) Info(format string, v ...any) {
	if LevelInformational > ll.level {
		return
	}
	msg := fmt.Sprintf("[I] "+format, v...)
	ll.Println(msg)
}

// Debug 校验
func (ll *Logger) Debug(format string, v ...any) {
	if LevelDebug > ll.level {
		return
	}
	msg := fmt.Sprintf("[D] "+format, v...)
	ll.Println(msg)
}

// BuildLogger 构建logger
func BuildLogger(loggerName string, level string) {
	intLevel := LevelError
	switch level {
	case "error":
		intLevel = LevelError
	case "warning":
		intLevel = LevelWarning
	case "info":
		intLevel = LevelInformational
	case "debug":
		intLevel = LevelDebug
	}
	l := Logger{
		name:  loggerName,
		level: intLevel,
	}
	logger = &l
}

// Log 返回日志对象
func Log() *Logger {
	if logger == nil {
		l := Logger{
			name:  "MediaCrawlerGo",
			level: LevelDebug,
		}
		logger = &l
	}
	return logger
}
