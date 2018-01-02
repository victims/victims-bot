package log

import "github.com/Sirupsen/logrus"

// Logger is the logger instance
var Logger *logrus.Logger

// init creates the Logger singleton
func init() {
	Logger = logrus.New()
	Logger.Level = logrus.DebugLevel
	Logger.Debug("Logger initialized")
}
