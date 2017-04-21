package logger

import (
	"fmt"
	"strings"
	"time"
)

const (
	STDOUT = iota
	STDERR
	PENTACONTA
)

var (
	LEVELS       = []string{"OUT", "ERR", "PEN"}
	LEVEL_COLORS = []string{"32", "31", "34"}
)

type Logger interface {
	Log(lg Log)
}

type Log struct {
	service string
	message string
	level   int
	time    time.Time
}

func NewLog(level int, service, message string) Log {
	return Log{
		time:    time.Now(),
		service: service,
		level:   level,
		message: message,
	}
}

type callLogger struct {
	logs     chan Log
	callback func(time.Time, string, string, string)
}

func NewCallLogger(callback func(time.Time, string, string, string)) Logger {
	l := &callLogger{
		logs:     make(chan Log, 100),
		callback: callback,
	}
	go l.start()
	return l
}

func (l *callLogger) start() {
	for lg := range l.logs {
		l.callback(lg.time, LEVELS[lg.level], lg.service, strings.Trim(lg.message, "\n"))
	}
}

func (l *callLogger) Log(lg Log) {
	l.logs <- lg
}

// stdoutLogger uses a channel to be able to log in correct order from different goroutines
// will print all logs to stdout
type stdoutLogger struct {
	logs chan Log
}

func NewStdoutLogger() Logger {
	l := &stdoutLogger{
		logs: make(chan Log, 100),
	}
	go l.start()
	return l
}

func (l *stdoutLogger) start() {
	for lg := range l.logs {
		timef := lg.time.Format("2006-01-02 15:04 05.999999")
		timef += strings.Repeat("0", 26-len(timef))

		cc := "\033[0;" + LEVEL_COLORS[lg.level] + "m"

		fmt.Printf("%v %v%v\033[0m %v: %v\n", timef, cc, LEVELS[lg.level], lg.service, strings.Trim(lg.message, "\n"))
	}
}

func (l *stdoutLogger) Log(lg Log) {
	l.logs <- lg
}
