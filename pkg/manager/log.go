package manager

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type (
	LogManager interface {
		Logrus() *logrus.Entry
		Echo() echo.Logger
	}

	LogOption interface {
		apply(l *logger)
	}

	logOption func(l *logger)

	logManager struct {
		l *logger
	}

	logger struct {
		*logrus.Entry
		prefix string
	}

	hook struct {
		writer io.Writer
		levels []logrus.Level
	}
)

func WithLogLevel(lvl logrus.Level) LogOption {
	return logOption(func(l *logger) {
		l.Logger.SetLevel(lvl)
	})
}

func WithLogFormatter(formatter logrus.Formatter) LogOption {
	return logOption(func(l *logger) {
		l.Logger.SetFormatter(formatter)
	})
}

func WithLogWriter(w io.Writer, levels ...logrus.Level) LogOption {
	return logOption(func(l *logger) {
		l.Logger.AddHook(&hook{
			writer: w,
			levels: levels,
		})
	})
}

func WithLogField(key string, value interface{}) LogOption {
	return logOption(func(l *logger) {
		l.Entry = l.Entry.WithField(key, value)
	})
}

func WithLogFields(fields logrus.Fields) LogOption {
	return logOption(func(l *logger) {
		l.Entry = l.Entry.WithFields(fields)
	})
}

func NewLogManager(opts ...LogOption) LogManager {
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

	lg := &logger{
		Entry:  l.WithField("hostname", hostname),
		prefix: "spring",
	}

	for _, opt := range opts {
		opt.apply(lg)
	}

	return &logManager{l: lg}
}

func (lm *logManager) Logrus() *logrus.Entry {
	return lm.l.Entry
}

func (lm *logManager) Echo() echo.Logger {
	return lm.l
}

func (fn logOption) apply(l *logger) {
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

func (l *logger) Output() io.Writer {
	return l.Entry.WriterLevel(l.Logger.GetLevel())
}

func (l *logger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

func (l *logger) Prefix() string {
	return l.prefix
}

func (l *logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *logger) Level() log.Lvl {
	return toEchoLevel(l.Logger.GetLevel())
}

func (l *logger) SetLevel(lvl log.Lvl) {
	l.Logger.SetLevel(toLogrusLevel(lvl))
}

func (l *logger) SetHeader(_ string) {
}

func (l *logger) Printj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Print()
}

func (l *logger) Debugj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Debug()
}

func (l *logger) Infoj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Info()
}

func (l *logger) Warnj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Warn()
}

func (l *logger) Errorj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Error()
}

func (l *logger) Fatalj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Fatal()
}

func (l *logger) Panicj(j log.JSON) {
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
