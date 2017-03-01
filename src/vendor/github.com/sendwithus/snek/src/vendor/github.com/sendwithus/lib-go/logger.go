package swu

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/getsentry/raven-go"
)

// Logger is commented
type Logger struct {
	slack                 *log.Logger
	info                  *log.Logger
	warn                  *log.Logger
	err                   *log.Logger
	fatal                 *log.Logger

	slackConfig           *SlackConfig
	slackService          slackService
	existingSlackChannels []string
}

// NewLogger returns a new instance of Logger
func NewLogger(sentryDsn string) *Logger {
	raven.SetDSN(sentryDsn)

	return &Logger{
		slackService:          NewSlackService(),
		existingSlackChannels: []string{},
		slack: log.New(os.Stdout, "[SLACK] ", log.Ldate|log.Ltime),
		info:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		warn:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime),
		err:   log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime),
		fatal: log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime),
	}
}

func (l *Logger) InitializeSlack(config *SlackConfig) error {
	l.slackConfig = config
	existingChannels, err := l.slackService.getSlackChannels(l.slackConfig.Token)
	if err != nil {
		return err
	}

	l.existingSlackChannels = existingChannels
	return nil
}

// Slack logs a message and sends to slack
func (l *Logger) Slack(channel string, format string, v ...interface{}) {
	if l.slackConfig == nil {
		msg := fmt.Sprintf(format, v)
		l.Error("could not log to slack channel [%v] message: %v", channel, msg)
		return
	}
	msg := fmt.Sprintf(format, v...)

	if !l.slackService.channelExists(channel, l.existingSlackChannels) {
		msg = fmt.Sprintf("Slack channel does not exist `%v`: %v", channel, msg)
		channel = "#fires"
	}

	l.slackService.sendMessage(
		l.slackConfig.Url,
		&SlackMessage{
			Text:     msg,
			Channel:  channel,
			Icon:     l.slackConfig.Icon,
			Username: l.slackConfig.Username,
		})

	l.slack.Printf(fmt.Sprintf("%s: %s", channel, msg))
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.warn.Printf(format, v...)
}

// Warn logs a warning message and an error
func (l *Logger) WarnWithError(err error, format string, v ...interface{}) {

	msg := fmt.Sprintf(format, v...)

	if msg == "" {
		msg = "Unknown"
	}

	l.warn.Printf("%s: %s", msg, err)
}

// Error logs an error
func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	err := errors.New(msg)

	// Log the error

	if msg == "" {
		msg = "Unknown"
	}

	l.err.Printf("%s", msg)

	raven.CaptureError(err, nil)
}

// ErrorWithError logs an error and captures it
func (l *Logger) ErrorWithError(err error, format string, v ...interface{}) {

	// Log the error

	msg := fmt.Sprintf(format, v...)

	if msg == "" {
		msg = "Unknown"
	}

	l.err.Printf("%s: %s", msg, err.Error())

	raven.CaptureError(err, nil)
}
