package telemetry

var Logger LoggerFace

type LoggerFace interface {
	Error(msg string)
	Errorf(format string, v ...interface{})
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Info(msg string)
	Infof(format string, v ...interface{})
	Debug(msg string)
	Debugf(format string, v ...interface{})
}
