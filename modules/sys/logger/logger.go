package logger

import (
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type (
	Option interface {
		apply(l *Logger)
	}
	
	Logger struct {
		*logrus.Entry
		prefix string
	}

	option func(l *Logger)

	hook struct {
		writer io.Writer
		levels []logrus.Level
	}
)

func WithDebug(debug bool) Option {
	lvl := logrus.InfoLevel
	if debug {
		lvl = logrus.DebugLevel
	}
	return WithLevel(lvl)
}

func WithLevel(lvl logrus.Level) Option {
	return option(func(l *Logger) {
		l.Logger.SetLevel(lvl)
	})
}

func WithFormatter(formatter logrus.Formatter) Option {
	return option(func(l *Logger) {
		l.Logger.SetFormatter(formatter)
	})
}

func WithWriter(w io.Writer, levels ...logrus.Level) Option {
	return option(func(l *Logger) {
		l.Logger.AddHook(&hook{
			writer: w,
			levels: levels,
		})
	})
}

func WithField(key string, value interface{}) Option {
	return option(func(l *Logger) {
		l.Entry = l.Entry.WithField(key, value)
	})
}

func WithFields(fields logrus.Fields) Option {
	return option(func(l *Logger) {
		l.Entry = l.Entry.WithFields(fields)
	})
}

func New(opts ...Option) *Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetOutput(ioutil.Discard)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	l.AddHook(&hook{
		writer: os.Stderr,
		levels: []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel},
	})

	l.AddHook(&hook{
		writer: os.Stdout,
		levels: []logrus.Level{logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel},
	})

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	logger := &Logger{
		Entry:  l.WithField("hostname", hostname),
		prefix: "spring",
	}

	for _, opt := range opts {
		opt.apply(logger)
	}

	return logger
}

func (fn option) apply(l *Logger) {
	fn(l)
}

func (h *hook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()

	if err != nil {
		return err
	}

	_, err = h.writer.Write(line)
	return err
}

func (h *hook) Levels() []logrus.Level {
	return h.levels
}

func (l *Logger) Output() io.Writer {
	return l.Entry.WriterLevel(l.Logger.GetLevel())
}

func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Level() log.Lvl {
	return toEchoLevel(l.Logger.GetLevel())
}

func (l *Logger) SetLevel(lvl log.Lvl) {
	l.Logger.SetLevel(toLogrusLevel(lvl))
}

func (l *Logger) SetHeader(_ string) {
}

func (l *Logger) Printj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Print()
}

func (l *Logger) Debugj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Debug()
}

func (l *Logger) Infoj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Info()
}

func (l *Logger) Warnj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Warn()
}

func (l *Logger) Errorj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Error()
}

func (l *Logger) Fatalj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Fatal()
}

func (l *Logger) Panicj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Panic()
}

func toLogrusLevel(level log.Lvl) logrus.Level {
	switch level {
	case log.DEBUG:
		return logrus.DebugLevel
	case log.INFO:
		return logrus.InfoLevel
	case log.WARN:
		return logrus.WarnLevel
	case log.ERROR:
		return logrus.ErrorLevel
	}

	return logrus.InfoLevel
}

func toEchoLevel(level logrus.Level) log.Lvl {
	switch level {
	case logrus.TraceLevel, logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return log.ERROR
	}

	return log.OFF
}
