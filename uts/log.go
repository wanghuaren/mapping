package uts

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

func InitLog(logName string) bool {
	f, err := os.OpenFile(logName+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return false
	}
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime)
	return true
}
func ChkErr(err error, args ...any) bool {
	if err != nil {
		_, code, line, ok := runtime.Caller(1)
		if !ok {
			LogF("获取runtime.Caller失败ChkErr")
		}
		errStr := getErrStr("错误", err.Error(), line, code, args...) + "\n"
		errStr += getErrStack()
		LogF(errStr)
		return true
	}
	return false
}

func ChkErrNormal(err error, args ...any) bool {
	if err != nil {
		_, code, line, ok := runtime.Caller(1)
		if !ok {
			LogF("获取runtime.Caller失败ChkErrNormal")
		}
		errStr := getErrStr("捕捉", err.Error(), line, code, args...)
		Log(errStr)
		return true
	}
	return false
}

func getErrStr(errTitle string, errDesc string, line int, code string, args ...any) string {
	errStr := ""
	errStr = errTitle + ":%v,行号%v,代码%v"
	errStr = fmt.Sprintf(errStr, errDesc, line, code)

	argsStr := ""
	for i := 0; i < len(args); i++ {
		argsStr += "%v,"
	}
	if argsStr != "" {
		errStr += "\n" + fmt.Sprintf(argsStr, args...)
	}
	return errStr
}

func ChkRecover() {
	_, code, line, ok := runtime.Caller(1)
	if !ok {
		LogF("获取runtime.Caller失败ChkRecover")
	}
	if !IsDebug {
		errDesc := ""
		if err := recover(); err != nil {
			errDesc = fmt.Sprintf("%v\n%v", err, string(debug.Stack()))
			panicStr := getErrStr("Panic", errDesc, line, code)
			panicStr += getErrStack()
			LogF(panicStr)
		}
	}
}

func getErrStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	stackStr := string(buf[:n])
	return stackStr
}

func LogF(args ...any) {
	LogDebug(args...)
}

func Log(args ...any) {
	logStr := LogDebug(args...)
	if !IsDebug {
		fmt.Println(logStr)
	}
}

func LogDebug(args ...any) string {
	logStr := ""
	if len(args) == 1 {
		logStr = "%v"
		logStr = fmt.Sprintf(logStr, args[0])
	} else if len(args) > 1 {
		if firstStr, ok := args[0].(string); ok && strings.Contains(firstStr, "%v") {
			paramPrint := []any{}
			for i := 1; i < len(args); i++ {
				paramPrint = append(paramPrint, args[i])
			}
			logStr = fmt.Sprintf(firstStr, paramPrint...)
		} else {
			for i := 0; i < len(args); i++ {
				logStr += "%v,"
			}
			logStr = logStr[:len(logStr)-1]
			logStr = fmt.Sprintf(logStr, args...)
		}

	}

	timeStr := time.Now().Format("2006-01-02 15:04:05.000")
	logStr = "[" + timeStr + "]" + logStr
	if IsDebug {
		fmt.Println(logStr)
	}
	log.Println(logStr)
	return logStr
}
