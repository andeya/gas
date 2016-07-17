package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	DEBUG = iota
	INFO
	WARNING
	DANGER
	ERROR
)

var (
	LOG_LEVEL = []string{
		"DEBUG",
		"INFO",
		"WARNING",
		"DANGER",
		"ERROR",
	}
)

type Logger struct {
	mu      *sync.Mutex
	logPath string
}

func New(log_path string) *Logger {
	l := &Logger{}
	l.mu = &sync.Mutex{}

	if log_path != "" {
		l.logPath = log_path
	} else {
		l.logPath = "/tmp/gas.log"
	}

	return l
}

func (l *Logger) LogF(level int, format string, data ...interface{}) {
	s := fmt.Sprintf(format, data...)

	l.doLog(level, s)
}

func (l *Logger) Debug(log string) {
	l.doLog(DEBUG, log)
}

func (l *Logger) Info(log string) {
	l.doLog(INFO, log)
}

func (l *Logger) Warning(log string) {
	l.doLog(WARNING, log)
}

func (l *Logger) Danger(log string) {
	l.doLog(DANGER, log)
}

func (l *Logger) Error(log string) {
	l.doLog(ERROR, log)
}

func (l *Logger) doLog(level int, log string) {
	go l.writeLog(level, log)
}

func (l *Logger) writeLog(level int, log string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// trim \n
	log = strings.Trim(log, "\n") + "\n"

	// add log mode
	log = fmt.Sprintf("[%s]%s", LOG_LEVEL[level], log)

	f, err := os.OpenFile(l.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 06666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err := f.WriteString(log); err != nil {
		fmt.Println("log error! Can not log to " + l.logPath + " file")

		return err
	}

	return nil
}
