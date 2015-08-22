// Package cologrus provides functionality to wrap Logrus hooks
// and formatters as ready to use CoLog hooks and formatters

package cologrus

import (
	"log"

	"comail.io/go/colog"
	"github.com/Sirupsen/logrus"
)

type logrusHook struct {
	hook   logrus.Hook
	levels []colog.Level
}

type logrusFormatter struct {
	formatter logrus.Formatter
}

// NewLogrusHook takes a Logrus hook and returns it wrapped as a CoLog hook
func NewLogrusHook(hook logrus.Hook) colog.Hook {
	levels := make([]colog.Level, len(hook.Levels()))
	for k, v := range hook.Levels() {
		levels[k] = level2colog(v)
	}

	return &logrusHook{hook: hook, levels: levels}
}

// Levels returns the levels for which this hook should be fired
func (lh *logrusHook) Levels() []colog.Level {
	return lh.levels
}

// Fire is called on the hook for every log entry matching the level
func (lh *logrusHook) Fire(entry *colog.Entry) error {
	return lh.hook.Fire(entry2logrus(entry))
}

// NewLogrusFormatter takes a Logrus formatter and returns it wrapped as a CoLog formatter
func NewLogrusFormatter(formatter logrus.Formatter) colog.Formatter {
	return &logrusFormatter{formatter: formatter}
}

func (lf *logrusFormatter) Format(entry *colog.Entry) ([]byte, error) {
	return lf.formatter.Format(entry2logrus(entry))
}

// SetFlags just implements the formatter interface
func (lf *logrusFormatter) SetFlags(flags int) {}

// Flags just implements the formatter interface
func (lf *logrusFormatter) Flags() int { return log.LstdFlags }

func level2colog(level logrus.Level) colog.Level {
	switch level {
	case logrus.DebugLevel:
		return colog.LDebug
	case logrus.InfoLevel:
		return colog.LInfo
	case logrus.WarnLevel:
		return colog.LWarning
	case logrus.ErrorLevel:
		return colog.LError
	}

	return colog.LAlert
}

func level2logrus(level colog.Level) logrus.Level {
	switch level {
	case colog.LTrace:
		return logrus.DebugLevel
	case colog.LDebug:
		return logrus.DebugLevel
	case colog.LInfo:
		return logrus.InfoLevel
	case colog.LWarning:
		return logrus.WarnLevel
	case colog.LError:
		return logrus.ErrorLevel
	}

	return logrus.FatalLevel
}

func entry2logrus(entry *colog.Entry) *logrus.Entry {
	data := make(logrus.Fields, len(entry.Fields))
	for k, v := range entry.Fields {
		data[k] = v
	}

	return &logrus.Entry{
		Logger:  nil,
		Data:    data,
		Time:    entry.Time,
		Level:   level2logrus(entry.Level),
		Message: string(entry.Message),
	}
}
