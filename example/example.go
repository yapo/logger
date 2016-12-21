package main

import (
	"github.schibsted.io/Yapo/logger"
	"time"
)

func main() {

	// syslog logs will be printed to /var/logs/messages by default
	conf := logger.LogConfig{logger.SyslogConfig{true, "exampleLogger"}, logger.StdlogConfig{true}}

	logger.Init(conf)
	logger.SetLogLevel(logger.INFO)

	logger.Debug("This is for Debugging")
	logger.Info("This is for Infoing")
	logger.Warn("This is a Warning")
	logger.Error("This is for Erroring")
	logger.Crit("This is for Critting")

	// this allows the example to print everything
	time.Sleep(1 * time.Second)
	logger.CloseSyslog()
}
