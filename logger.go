package fk_logrus_formatter

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Constants for FkLogFormatter
const (
	logrusStackJump          = 4
	logrusFieldlessStackJump = 6
	FunctionKey              = "function"
	PackageKey               = "package"
	LineKey                  = "line"
	FileKey                  = "file"
)

type level int

// Constants for Levels
const (
	PANIC level = iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	// In FkLogFormatter you can specify
	runtimeFormatter := formatter{ChildFormatter: FkLogFormatter{}}
	log.Formatter = &runtimeFormatter
}

// GetLogger returns a copy of the log with the given level
func GetLogger(level logrus.Level) logrus.Logger {
	lCopy := *log
	lCopy.SetLevel(level)
	return lCopy
}

func SetFileLogging(logFile string) error {
	err := os.MkdirAll(path.Dir(logFile), os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create log directory. Err %v", err)
		return err
	}

	logFileP, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(logFileP)
	log.Infof("Logging started")
	return nil
}

// Close closes the log file handle
func Close() {
	log.Infof("Closing log file")
	fp, ok := log.Out.(*os.File)
	if ok {
		fp.Close()
	}
}

func SetLevel(ll level) {
	switch ll {
	case PANIC:
		log.SetLevel(logrus.InfoLevel)
	case FATAL:
		log.SetLevel(logrus.FatalLevel)
	case ERROR:
		log.SetLevel(logrus.ErrorLevel)
	case WARN:
		log.SetLevel(logrus.WarnLevel)
	case INFO:
		log.SetLevel(logrus.InfoLevel)
	case DEBUG:
		log.SetLevel(logrus.DebugLevel)
	}
}

// Infof INFO level log printing
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf WARN level log printing
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf ERROR level log printing
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Debugf DEBUG level log printing
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func IsDebugEnabled() bool {
	return log.IsLevelEnabled(logrus.DebugLevel)
}

// Debugf PANIC level log printing
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Fatalf FATAL level log printing
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Adapted from https://github.com/banzaicloud/logrus-runtime-formatter
type formatter struct {
	ChildFormatter FkLogFormatter
	Line           bool
	Package        bool
	File           bool
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// you can use any other specific function likewise getCurrentPosition to get any other details
	function, file, line := f.getCurrentPosition(entry)

	packageEnd := strings.LastIndex(function, ".")
	functionName := function[packageEnd+1:]

	data := logrus.Fields{FunctionKey: functionName}
	if f.Line {
		data[LineKey] = line
	}
	if f.Package {
		packageName := function[:packageEnd]
		data[PackageKey] = packageName
	}
	if f.File {
		pathTokens := strings.Split(file, "/")
		data[FileKey] = pathTokens[len(pathTokens)-1]
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	entry.Data = data

	return f.ChildFormatter.Format(entry)
}

func (f *formatter) getCurrentPosition(entry *logrus.Entry) (string, string, string) {
	// skip = skip + 2, +2 is for one more level of indirection added in Altair
	skip := logrusStackJump + 2
	if len(entry.Data) == 0 {
		skip = logrusFieldlessStackJump + 2
	}
start:
	pc, file, line, _ := runtime.Caller(skip)
	lineNumber := ""
	if f.Line {
		lineNumber = fmt.Sprintf("%d", line)
	}
	function := runtime.FuncForPC(pc).Name()
	if strings.LastIndex(function, "sirupsen/logrus.") != -1 {
		skip++
		goto start
	}
	return function, file, lineNumber
}
