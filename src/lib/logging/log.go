package logging

import (
	"log"
	"os"

	"github.com/epsniff/gui/src/lib/logging/loggergou"
)

var (
	verbose *bool
)

//Logger used by all other packages.
// Defaults to a gou logger that only logs fatal logs effectively disableing it,
// but we need something here to avoid panics.
// TODO default to a NOOP logger.
var Logger LoggerFace = loggergou.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds), "fatal")

func init() {
	if os.Getenv("VERBOSELOGS") != "" {
		//mostly this ENV is here for turning on logging in testcases.
		Logger = loggergou.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds), "debug")
	}
}

type LoggerFace interface {
	Error(msg string)
	Errorf(format string, v ...interface{})
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Info(msg string)
	Infof(format string, v ...interface{})
	Debug(msg string)
	Debugf(format string, v ...interface{})
	Printf(format string, v ...interface{})
}
