package logbeat

import (
	"github.com/sirupsen/logrus"
)

// LogbeatVersion is used to identify notifications sent by Logbeat.
const LogbeatVersion string = "0.0.3"

// LogbeatHook delivers logs to the Opbeat service.
type LogbeatHook struct {
	AppId       string
	Opbeat      *OpbeatClient
	OrgId       string
	SecretToken string
}

// NewOpbeatHook creates a new LogbeatHook that reports
// Errors, Fatal Errors, and Panics from Logrus to Opbeat.
func NewOpbeatHook(org, app, token string) *LogbeatHook {
	return &LogbeatHook{
		OrgId:       org,
		AppId:       app,
		SecretToken: token,
		Opbeat:      NewOpbeatClient(org, app, token),
	}
}

// Fire is called by Logrus when an Error occurs. The given Logrus Entry
// will be sent to Opbeat's Error Intake API.
func (hook *LogbeatHook) Fire(entry *logrus.Entry) error {
	_, err := hook.Opbeat.Notify(entry)
	return err
}

// Levels tells Logrus which types of logging events we
// are interseted in (e.g. Errors, Fatal Errors, Panics).
func (hook *LogbeatHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}
