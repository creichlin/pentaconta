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
	LEVELS      = []string{"OUT", "ERR", "PEN"}
	ANSI_LEVELS = map[string]string{
		"OUT": "\033[0;32OUT\033[0m",
		"ERR": "\033[0;31ERR\033[0m",
		"PEN": "\033[0;34PEN\033[0m",
	}
)

type Logger interface {
	Log(lg Log)
}

type Log struct {
	Service  string
	Instance int
	Message  string
	Level    int
	Time     time.Time
}

func NewLog(level int, service string, instance int, message string) Log {
	return Log{
		Time:     time.Now(),
		Service:  service,
		Instance: instance,
		Level:    level,
		Message:  message,
	}
}

type splitLogger struct {
	targets []Logger
}

func NewSplitLogger(targets ...Logger) Logger {
	return &splitLogger{
		targets: targets,
	}
}

func (l *splitLogger) Log(lg Log) {
	for _, l := range l.targets {
		l.Log(lg)
	}
}

type callLogger struct {
	logs     chan Log
	callback func(time.Time, string, string, int, string)
}

func NewCallLogger(callback func(time.Time, string, string, int, string)) Logger {
	l := &callLogger{
		logs:     make(chan Log, 100),
		callback: callback,
	}
	go l.start()
	return l
}

func (l *callLogger) start() {
	for lg := range l.logs {
		l.callback(lg.Time, LEVELS[lg.Level], lg.Service, lg.Instance, strings.Trim(lg.Message, "\n"))
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
		timef := lg.Time.Format("2006-01-02 15:04 05.999999")
		timef += strings.Repeat("0", 26-len(timef))
		fmt.Printf("%v %v %v%v: %v\n", timef, ANSI_LEVELS[LEVELS[lg.Level]], lg.Service, lg.Instance, strings.Trim(lg.Message, "\n"))
	}
}

func (l *stdoutLogger) Log(lg Log) {
	l.logs <- lg
}
