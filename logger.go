package logger

import (
	"fmt"
	"log"
	"log/syslog"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Conf structure for syslog options
// Indentity is the syslog tag
type SyslogConfig struct {
	Enabled  bool
	Identity string
}

// Conf structure for stdout logging
type StdlogConfig struct {
	Enabled bool
}

// General conf structure
type LogConfig struct {
	Syslog SyslogConfig
	Stdlog StdlogConfig
}

// Constant with the different log levels
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	CRIT
	NUMOFLOGLEVELS
)

type logType struct {
	Channel      chan string
	SyslogHandle func(msg string) error
	StdlogHandle func(msg string)
}

var (
	logTypes      [NUMOFLOGLEVELS]logType
	logSys        *syslog.Writer
	loggerRunning = false
	loggerLevel   = DEBUG
)

// Set the log level
func SetLogLevel(level int) {
	loggerLevel = level
}

func logAtLevel(debugLevel int, format string, params ...interface{}) {
	format = funcName() + format
	if loggerRunning && debugLevel >= loggerLevel {
		if len(params) > 0 {
			logTypes[debugLevel].Channel <- fmt.Sprintf(format, params...)
		} else {
			logTypes[debugLevel].Channel <- format
		}
	}
}

// Logs a critical message
func Crit(format string, params ...interface{}) {
	logAtLevel(CRIT, format, params...)
}

// Logs debuging messages
func Debug(format string, params ...interface{}) {
	logAtLevel(DEBUG, format, params...)
}

// Logs an Error message
func Error(format string, params ...interface{}) {
	logAtLevel(ERROR, format, params...)
}

// Logs an Info message
func Info(format string, params ...interface{}) {
	logAtLevel(INFO, format, params...)
}

// logs a warning message
func Warn(format string, params ...interface{}) {
	logAtLevel(WARN, format, params...)
}

// Inits the logger module
func Init(conf LogConfig) error {
	if conf.Syslog.Enabled {
		var err error
		logSys, err = syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, conf.Syslog.Identity)
		if err != nil {
			return fmt.Errorf("Logger - Error: %s", err)
		}

		logTypes[CRIT].SyslogHandle = logSys.Crit
		logTypes[DEBUG].SyslogHandle = logSys.Debug
		logTypes[ERROR].SyslogHandle = logSys.Err
		logTypes[INFO].SyslogHandle = logSys.Info
		logTypes[WARN].SyslogHandle = logSys.Warning

		loggerRunning = true
	}

	if conf.Stdlog.Enabled {
		logTypes[CRIT].StdlogHandle = logCrit
		logTypes[DEBUG].StdlogHandle = logDebug
		logTypes[ERROR].StdlogHandle = logError
		logTypes[INFO].StdlogHandle = logInfo
		logTypes[WARN].StdlogHandle = logWarn

		loggerRunning = true
	}

	if !loggerRunning {
		return fmt.Errorf("Logger - Error: Not running")
	}

	for i := 0; i < NUMOFLOGLEVELS; i++ {
		logTypes[i].Channel = make(chan string)
	}
	go consumeLogs(conf)
	return nil
}

// close the Syslog connection
func CloseSyslog() error {
	time.Sleep(time.Second)
	err := logSys.Close()
	if err != nil {
		return fmt.Errorf("Logger - Error: %s", err)
	}
	return err
}

func consumeLogs(conf LogConfig) {

	logToSubsystems := func(level int, msg string) {
		if conf.Syslog.Enabled {
			err := logTypes[level].SyslogHandle(msg)
			if err != nil {
				fmt.Printf(err.Error())
			}
		}
		if conf.Stdlog.Enabled {
			logTypes[level].StdlogHandle(msg)
		}
	}

	for {
		select {
		case msg := <-logTypes[CRIT].Channel:
			logToSubsystems(CRIT, msg)

		case msg := <-logTypes[DEBUG].Channel:
			logToSubsystems(DEBUG, msg)

		case msg := <-logTypes[ERROR].Channel:
			logToSubsystems(ERROR, msg)

		case msg := <-logTypes[INFO].Channel:
			logToSubsystems(INFO, msg)

		case msg := <-logTypes[WARN].Channel:
			logToSubsystems(WARN, msg)

		}
	}
}

func logCrit(msg string) {
	log.Printf("CRIT: %s\n", msg)
}

func logDebug(msg string) {
	log.Printf("DEBUG: %s\n", msg)
}

func logError(msg string) {
	log.Printf("ERROR: %s\n", msg)
}

func logInfo(msg string) {
	log.Printf("INFO: %s\n", msg)
}

func logWarn(msg string) {
	log.Printf("WARN: %s\n", msg)
}

func funcName() string {
	pc, _, _, ok := runtime.Caller(3)
	if ok {
		funcPtr := runtime.FuncForPC(pc)
		if funcPtr != nil {
			nameEnd := filepath.Ext(funcPtr.Name())
			return strings.TrimPrefix(nameEnd, ".") + ": "
		}
	}
	return ""
}
