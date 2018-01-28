package loggergou

import (
	"fmt"
	"log"

	"github.com/araddon/gou"
)

func New(l *log.Logger, logLevel string) *GouLogger {
	//gou.SetLogger(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds), "info")
	gou.SetLogger(l, logLevel)
	gou.SetColorOutput()

	return &GouLogger{}
}

type GouLogger struct{}

const linelvl = 3

func (l *GouLogger) Error(msg string) {
	if gou.LogLevel >= gou.ERROR {
		gou.DoLog(linelvl, gou.ERROR, msg)
	}
}
func (l *GouLogger) Errorf(format string, v ...interface{}) {
	if gou.LogLevel >= gou.ERROR {
		gou.DoLog(linelvl, gou.ERROR, fmt.Sprintf(format, v...))
	}
}

func (l *GouLogger) Warn(msg string) {
	if gou.LogLevel >= gou.WARN {
		gou.DoLog(linelvl, gou.WARN, msg)
	}
}
func (l *GouLogger) Warnf(format string, v ...interface{}) {
	if gou.LogLevel >= gou.WARN {
		gou.DoLog(linelvl, gou.WARN, fmt.Sprintf(format, v...))
	}
}

func (l *GouLogger) Info(msg string) {
	if gou.LogLevel >= gou.INFO {
		gou.DoLog(linelvl, gou.INFO, msg)
	}
}
func (l *GouLogger) Infof(format string, v ...interface{}) {
	if gou.LogLevel >= gou.INFO {
		gou.DoLog(linelvl, gou.INFO, fmt.Sprintf(format, v...))
	}
}

func (l *GouLogger) Debug(msg string) {
	if gou.LogLevel >= gou.DEBUG {
		gou.DoLog(linelvl, gou.DEBUG, msg)
	}
}
func (l *GouLogger) Debugf(format string, v ...interface{}) {
	if gou.LogLevel >= gou.DEBUG {
		gou.DoLog(linelvl, gou.DEBUG, fmt.Sprintf(format, v...))
	}
}
