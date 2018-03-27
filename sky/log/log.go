package log

import (
	"io"
	"os"
	"fmt"
	"strings"
	ml "log"
)

const (
	DEBUG = iota
	TRACE
	WARNING
	ERROR
	CRITICAL
)


type Logger struct {
	*ml.Logger
	Lv 	int
}

var	levels = [...]string{"DEBUG", "TRACE", "WARN", "ERROR", "CRITICAL"}
var	outs   = [...]string{"[DEBG] ", "[TRAC] ", "[WARN] ", "[EROR] ", "[CRIT] "}

var global *Logger = NewLogger(os.Stdout, "debug", 19)

func NewLogger(w io.Writer, lv string, flag int) *Logger {
	return &Logger{ml.New(w, "", flag), stringLevel(lv),}
}

func NewGlobal(w io.Writer, lv string, flag int) {
	global =  &Logger{ml.New(w, "", flag), stringLevel(lv),}
}

func SetGlobalLevel(lv string) {
	global.Lv = stringLevel(lv)
}

func Debug(format interface{}, v ...interface{}) {
	if global.Lv > DEBUG { return }
	global.Logger.SetPrefix(outs[DEBUG])
	switch format := format.(type) {
	case string:
		global.Logger.Output(2, fmt.Sprintf(format, v...))
	default:
		global.Logger.Output(2, fmt.Sprintf("%s", format))
	}	
}

func Trace(format interface{}, v ...interface{}) {
	if global.Lv > TRACE { return }
	global.Logger.SetPrefix(outs[TRACE])
	switch format := format.(type) {
	case string:
		global.Logger.Output(2, fmt.Sprintf(format, v...))
	default:
		global.Logger.Output(2, fmt.Sprintf("%s", format))
	}	
}

func Warn(format interface{}, v ...interface{}) {
	if global.Lv > WARNING { return }
	global.Logger.SetPrefix(outs[WARNING])
	switch format := format.(type) {
	case string:
		global.Logger.Output(2, fmt.Sprintf(format, v...))
	default:
		global.Logger.Output(2, fmt.Sprintf("%s", format))
	}	
}

func Error(format interface{}, v ...interface{}) {
	if global.Lv > ERROR { return }
	global.Logger.SetPrefix(outs[ERROR])
	switch format := format.(type) {
	case string:
		global.Logger.Output(2, fmt.Sprintf(format, v...))
	default:
		global.Logger.Output(2, fmt.Sprintf("%s", format))
	}		
}

func Critical(format interface{}, v ...interface{}) {
	global.Logger.SetPrefix(outs[CRITICAL])
	switch format := format.(type) {
	case string:
		global.Logger.Output(2, fmt.Sprintf(format, v...))
	default:
		global.Logger.Output(2, fmt.Sprintf("%s", format))
	}		
}

func SetLevel(lv int) {
	global.Lv = lv
}

func levelString(l int) string {
	if l < DEBUG || l > CRITICAL {
		return "UNKNOWN"
	}
	return levels[l]
}

func stringLevel(s string) int {
	ss := strings.ToUpper(s)
	for i := range levels {
		if ss == levels[i] { 
			return i
		}		
	}
	return DEBUG
}
