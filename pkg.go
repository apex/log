package log

// singletons ftw?
var std = &Logger{
	Level: InfoLevel,
}

// SetHandler sets the handler. This is not thread-safe.
func SetHandler(h Handler) {
	std.Handler = h
}

// SetLevel sets the log level. This is not thread-safe.
func SetLevel(l Level) {
	std.Level = l
}

// WithFields returns a new entry with `fields` set.
func WithFields(fields Fielder) *Entry {
	return std.WithFields(fields)
}

// WithField returns a new entry with the `key` and `value` set.
func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

// WithError returns a new entry with the "error" set to `err`.
func WithError(err error) *Entry {
	return std.WithError(err)
}

// Debug level message.
func Debug(msg string) {
	std.Debug(msg)
}

// Info level message.
func Info(msg string) {
	std.Info(msg)
}

// Warn level message.
func Warn(msg string) {
	std.Warn(msg)
}

// Error level message.
func Error(msg string) {
	std.Error(msg)
}

// Fatal level message, followed by an exit.
func Fatal(msg string) {
	std.Fatal(msg)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	std.Debugf(msg, v...)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	std.Infof(msg, v...)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	std.Warnf(msg, v...)
}

// Errorf level formatted message.
func Errorf(msg string, v ...interface{}) {
	std.Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	std.Fatalf(msg, v...)
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func Trace(msg string) *Entry {
	return std.Trace(msg)
}
