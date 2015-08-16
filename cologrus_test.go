package cologrus

import (
	"testing"
	"time"

	"strings"

	"comail.io/go/colog"
	"github.com/Sirupsen/logrus"
)

type LogrusInterfaces struct {
	formatter logrus.Formatter
	hook      logrus.Hook
}

type LogrusFormatter struct{}

func (l *LogrusFormatter) Format(e *logrus.Entry) ([]byte, error) { return nil, nil }

type LogrusHook struct {
	Entry *logrus.Entry
	Lvls  []logrus.Level
}

func (l *LogrusHook) Levels() []logrus.Level     { return l.Lvls }
func (l *LogrusHook) Fire(e *logrus.Entry) error { l.Entry = e; return nil }

// Test interfaces are met
var i = LogrusInterfaces{formatter: new(LogrusFormatter), hook: new(LogrusHook)}

// TTime is the fixed point in time for all formatting tests
var TTime = time.Date(2015, time.August, 1, 20, 45, 30, 9999, time.UTC)

func TestFormatter(t *testing.T) {
	e := &colog.Entry{
		Time:    TTime,
		Level:   colog.LDebug,
		Message: []byte("some message"),
		Fields:  map[string]interface{}{"foo": "bar"},
	}

	cf := NewLogrusFormatter(&logrus.TextFormatter{DisableColors: true})
	data, err := cf.Format(e)
	if err != nil {
		t.Errorf("Error formatting entry %s", err.Error())
	}

	dataStr := strings.Trim(string(data), " \n")
	if dataStr != `time="2015-08-01T20:45:30Z" level=debug msg="some message" foo=bar` {
		t.Errorf("Invalid formatter entry '%s'", dataStr)
	}

	// Double check that this is all still compatible
	cf = NewLogrusFormatter(new(LogrusFormatter))
	logrus.SetFormatter(new(LogrusFormatter))
	colog.SetFormatter(cf)
}

func TestHook(t *testing.T) {
	e := &colog.Entry{
		Time:    TTime,
		Level:   colog.LDebug,
		Message: []byte("some message"),
		Fields:  map[string]interface{}{"foo": "bar"},
	}

	lh := &LogrusHook{Lvls: []logrus.Level{
		logrus.DebugLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}}

	ch := NewLogrusHook(lh)

	if len(ch.Levels()) != 3 {
		t.Errorf("Missing levels, found %d", len(ch.Levels()))
	}

	err := ch.Fire(e)
	if err != nil {
		t.Errorf("Error firing hook %s", err.Error())
	}

	if string(lh.Entry.Message) != string(e.Message) {
		t.Errorf("Invalid message on entry %s", lh.Entry.Message)
	}

	if lh.Entry.Data["foo"].(string) != "bar" {
		t.Errorf("Invalid foo field on entry %s", lh.Entry.Data)
	}

	if lh.Entry.Data["foo"].(string) != "bar" {
		t.Errorf("Invalid foo field on entry %s", lh.Entry.Data)
	}

	if lh.Entry.Time != e.Time {
		t.Errorf("Invalid time on entry %s", lh.Entry.Time)
	}

	if lh.Entry.Level != logrus.DebugLevel {
		t.Errorf("Invalid level on entry %s", lh.Entry.Level)
	}

	// Double check that this is all still compatible
	logrus.AddHook(lh)
	colog.AddHook(ch)
}
