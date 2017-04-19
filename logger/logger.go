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

type Log struct {
	service string
	message string
	level   int
	time    time.Time
}

// Logger uses a channel to be able to lof in correct order from different goroutines
type Logger struct {
	logs chan Log
}

func NewLogger() *Logger {
	l := &Logger{
		logs: make(chan Log, 100),
	}
	go l.start()
	return l
}

func (l *Logger) start() {
	for lg := range l.logs {
		timef := lg.time.Format("2006-01-02 15:04 05.999999")
		timef += strings.Repeat("0", 26-len(timef))

		cc := "\033[0;" + LEVEL_COLORS[lg.level] + "m"

		fmt.Printf("%v %v%v\033[0m %v: %v\n", timef, cc, LEVELS[lg.level], lg.service, strings.Trim(lg.message, "\n"))
	}
}

func NewLog(level int, service, message string) Log {
	return Log{
		time:    time.Now(),
		service: service,
		level:   level,
		message: message,
	}
}

func (l *Logger) Log(lg Log) {
	l.logs <- lg
}
