package logger

import (
	"fmt"
	"log"
)

type LoggerI interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

var levels = map[Level]level{
	debug: debugLevel,
	info:  infoLevel,
	warn:  warnLevel,
	err:   errLevel,
	fatal: fatalLevel,
}

type Level string
type level int

const (
	warnLevel level = iota
	errLevel
	fatalLevel
	infoLevel
	debugLevel
)

const (
	debug Level = "debug"
	info  Level = "info"
	warn  Level = "warn"
	err   Level = "error"
	fatal Level = "fatal"
)

type Logger struct {
	Level  Level
	Prefix string
	l      *log.Logger
}

func New(l string, p string) *Logger {
	return &Logger{Level: Level(l), Prefix: p, l: log.Default()}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if levels[l.Level] < debugLevel {
		return
	}
	l.l.SetPrefix(fmt.Sprintf("[%s] ", debug))
	if args != nil {
		msg = fmt.Sprintf(msg, args)
	}
	msg = fmt.Sprintf("%s\n", msg)
	l.l.Printf(msg)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if levels[l.Level] < infoLevel {
		return
	}
	l.l.SetPrefix(fmt.Sprintf("[%s] ", info))
	if args != nil {
		msg = fmt.Sprintf(msg, args)
	}
	msg = fmt.Sprintf("%s\n", msg)
	l.l.Printf(msg)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if levels[l.Level] < warnLevel {
		return
	}
	l.l.SetPrefix(fmt.Sprintf("[%s] ", warn))
	if args != nil {
		msg = fmt.Sprintf(msg, args)
	}
	msg = fmt.Sprintf("%s\n", msg)
	l.l.Printf(msg)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if levels[l.Level] < errLevel {
		return
	}
	l.l.SetPrefix(fmt.Sprintf("[%s] ", err))
	if args != nil {
		msg = fmt.Sprintf(msg, args)
	}
	msg = fmt.Sprintf("%s\n", msg)
	l.l.Printf(msg)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.l.SetPrefix(fmt.Sprintf("[%s] ", fatal))
	if args != nil {
		msg = fmt.Sprintf(msg, args)
	}
	msg = fmt.Sprintf("%s\n", msg)
	l.l.Fatalf(msg)
}
