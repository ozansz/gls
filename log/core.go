package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

var (
	out   io.Writer = os.Stdout
	debug int32     = 0
)

func SetDebug(d int32) {
	atomic.StoreInt32(&debug, d)
}

func SetOutput(w io.Writer) {
	out = w
}

func Debug(args ...interface{}) {
	if atomic.LoadInt32(&debug) != 0 {
		fmt.Fprintln(out, fmt.Sprintf("[DEBUG] "+fmt.Sprint(args...)))
	}
}

func Debugf(format string, args ...interface{}) {
	if atomic.LoadInt32(&debug) != 0 {
		fmt.Fprintln(out, fmt.Sprintf("[DEBUG] "+format, args...))
	}
}

func Info(args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[INFO] "+fmt.Sprint(args...)))
}

func Infof(format string, args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[INFO] "+format, args...))
}

func Warning(args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[WARNING] "+fmt.Sprint(args...)))
}

func Warningf(format string, args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[WARNING] "+format, args...))
}

func Error(args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[ERROR] "+fmt.Sprint(args...)))
}

func Errorf(format string, args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[ERROR] "+format, args...))
}

func Fatal(args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[FATAL] "+fmt.Sprint(args...)))
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Fprintln(out, fmt.Sprintf("[FATAL] "+format, args...))
	os.Exit(1)
}
