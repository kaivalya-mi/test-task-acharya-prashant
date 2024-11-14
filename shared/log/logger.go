package log

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Data is
type Data struct {
	IPAddress string `` // the ip address of caller
	Session   string `` // id that generated in controller and passed to service to service for flow tracking purpose
	ActorID   string `` // could be userID from Apps or backoffice
	ActorType string `` // MOB (Mobile Apps) / BOF (Backoffice) / MSQ (message queuing) / SCH (scheduller) / SYS (system)
}

// ILogger is
type ILogger interface {
	Debug(data interface{}, description string, args ...interface{})
	Info(data interface{}, description string, args ...interface{})
	Warn(data interface{}, description string, args ...interface{})
	Error(data interface{}, description string, args ...interface{})
	Fatal(data interface{}, description string, args ...interface{})
	Panic(data interface{}, description string, args ...interface{})
}

// LogrusImpl is
type logger struct {
	appName    string
	appVersion string

	filePath string
	level    logrus.Level
	maxAge   time.Duration

	theLogger *logrus.Logger
}

var defaultLogger logger
var defaultLoggerOnce sync.Once

func createLogger() {
	// formatter := logrus.JSONFormatter{}

	formatter := nested.Formatter{
		NoColors:        true,
		HideKeys:        true,
		TimestampFormat: "0102 150405.000",
		FieldsOrder:     []string{"func"},
	}

	defaultLogger.theLogger = logrus.New()
	defaultLogger.theLogger.SetLevel(defaultLogger.level)
	defaultLogger.theLogger.SetFormatter(&formatter)

	filename := defaultLogger.appName + ".%Y%m%d.log"
	if len(defaultLogger.appVersion) > 0 {
		filename = defaultLogger.appName + "-" + defaultLogger.appVersion + ".%Y%m%d.log"
	}

	writer, _ := rotatelogs.New(
		filepath.Join(defaultLogger.filePath, filename),
		// rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(defaultLogger.maxAge),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	defaultLogger.theLogger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.DebugLevel: writer,
		},
		defaultLogger.theLogger.Formatter,
	))
}

// WithFile is command to state the log will printing to files
// the rolling log file will put in logs/ directory
//
// filename is just a name of log file without any extension
//
// maxAge is age (in days) of the logs file before it gets purged from the file system
func Init(appName, appVersion, filePath string, logLevel logrus.Level, maxAge time.Duration) {
	defaultLogger.appName = appName
	defaultLogger.appVersion = appVersion
	defaultLogger.filePath = filePath
	defaultLogger.level = logLevel
	defaultLogger.maxAge = maxAge
}

// GetLog is
func GetLog() ILogger {
	defaultLoggerOnce.Do(createLogger)
	return &defaultLogger
}

func (self *logger) getLogEntry(extraInfo interface{}) *logrus.Entry {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	var buffer bytes.Buffer

	buffer.WriteString("fn:")

	x := strings.LastIndex(funcName, "/")
	buffer.WriteString(funcName[x+1:])

	if extraInfo == nil {
		return self.theLogger.WithField("info", buffer.String())
	}

	data, ok := extraInfo.(Data)
	if !ok {
		return self.theLogger.WithField("info", buffer.String())
	}

	if data.IPAddress != "" {
		buffer.WriteString("|ip:")
		buffer.WriteString(data.IPAddress)
	}

	if data.Session != "" {
		buffer.WriteString("|ss:")
		buffer.WriteString(data.Session)
	}

	if data.ActorID != "" {
		buffer.WriteString("|id:")
		buffer.WriteString(data.ActorID)
	}

	if data.ActorType != "" {
		buffer.WriteString("|tp:")
		buffer.WriteString(data.ActorType)
	}

	return self.theLogger.WithField("info", buffer.String())
}

// Debug is
func (self *logger) Debug(data interface{}, description string, args ...interface{}) {
	self.getLogEntry(data).Debugf(description+"\n", args...)
}

// Info is
func (self *logger) Info(data interface{}, description string, args ...interface{}) {
	self.getLogEntry(data).Infof(description+"\n", args...)
}

// Warn is
func (self *logger) Warn(data interface{}, description string, args ...interface{}) {
	self.getLogEntry(data).Warnf(description+"\n", args...)
}

// Error is
func (self *logger) Error(data interface{}, description string, args ...interface{}) {
	self.getLogEntry(data).Errorf(description+"\n", args...)
}

// Fatal is
func (self *logger) Fatal(data interface{}, description string, args ...interface{}) {
	self.getLogEntry(data).Fatalf(description+"\n", args...)
}

// Panic is
func (self *logger) Panic(data interface{}, description string, args ...interface{}) {
	self.getLogEntry(data).Panicf(description+"\n", args...)
}
