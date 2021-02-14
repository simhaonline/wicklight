package logger

import (
	"io"
	"log"
	"os"
)

// LogLevel is the level of logger
type LogLevel int

// Levels pre-defined
const (
	DebugLevel LogLevel = 1
	InfoLevel  LogLevel = 2
	WarnLevel  LogLevel = 3
	ErrorLevel LogLevel = 4
	FatalLevel LogLevel = 5
	MuteLevel  LogLevel = 6
)

var (
	currentLevel  LogLevel  = 2
	currentOutput io.Writer = os.Stdout
)

func init() {
	log.SetFlags(log.LstdFlags)
	SetLevel(InfoLevel)
	SetOutput("")
}

// SetLevel sets the current log level
func SetLevel(level LogLevel) {
	currentLevel = level
}

// SetOutput sets the output file for logger
func SetOutput(filename string) {
	if filename == "" {
		currentOutput = os.Stdout
	} else {
		f, err := os.Create(filename)
		if err != nil {
			Fatal("[log] can not open log file:", filename)
		}
		currentOutput = f
	}
	log.SetOutput(currentOutput)
}

// Fatal fatal logging
func Fatal(v ...interface{}) {
	if currentLevel <= FatalLevel {
		log.Fatalln(v...)
	}
}

// Fatalf fatal logging
func Fatalf(format string, v ...interface{}) {
	if currentLevel <= FatalLevel {
		log.Fatalf(format, v...)
	}
}

// Error error logging
func Error(v ...interface{}) {
	if currentLevel <= ErrorLevel {
		log.Println(v...)
	}
}

// Errorf error logging
func Errorf(format string, v ...interface{}) {
	if currentLevel <= ErrorLevel {
		log.Printf(format, v...)
	}
}

// Warn warn logging
func Warn(v ...interface{}) {
	if currentLevel <= WarnLevel {
		log.Println(v...)
	}
}

// Warnf warn logging
func Warnf(format string, v ...interface{}) {
	if currentLevel <= WarnLevel {
		log.Printf(format, v...)
	}
}

// Info info logging
func Info(v ...interface{}) {
	if currentLevel <= InfoLevel {
		log.Println(v...)
	}
}

// Infof info logging
func Infof(format string, v ...interface{}) {
	if currentLevel <= InfoLevel {
		log.Printf(format, v...)
	}
}

// Debug debug logging
func Debug(v ...interface{}) {
	if currentLevel <= DebugLevel {
		log.Println(v...)
	}
}

// Debugf debug logging
func Debugf(format string, v ...interface{}) {
	if currentLevel <= DebugLevel {
		log.Printf(format, v...)
	}
}
