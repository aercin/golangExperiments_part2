package logrus

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Log(msg string, cfg *Configs) error
}

type logger struct {
	log *logrus.Logger
}

func NewLogger(minLevel level, hooks ...any) (Logger, error) {

	entry := logrus.New()

	entry.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	entry.SetLevel(logrus.Level(minLevel))

	for _, hook := range hooks {
		if h, ok := hook.(logrus.Hook); ok {
			entry.AddHook(h)
		} else {
			return nil, fmt.Errorf("provided hook does not implement logrus.Hook interface")
		}
	}
	// if hook != nil {
	// 	if hook, ok := hook.(logrus.Hook); ok {
	// 		entry.AddHook(hook)
	// 	} else {
	// 		return nil, fmt.Errorf("provided hook does not implement logrus.Hook interface")
	// 	}
	// }

	return &logger{
		log: entry,
	}, nil
}

func (l *logger) Log(msg string, cfg *Configs) error {

	logFields := make(logrus.Fields)

	for key, value := range cfg.CustomEntries {
		logFields[key] = value
	}

	entry := l.log.WithFields(logFields)

	switch cfg.Level {
	case Trace:
		entry.Trace(msg)
	case Debug:
		entry.Debug(msg)
	case Info:
		entry.Info(msg)
	case Warn:
		entry.Warn(msg)
	case Error:
		entry.Error(msg)
	case Fatal:
		entry.Fatal(msg)
	default:
		return errors.New("An error occured: passing undefined level parameter value")
	}

	return nil
}
