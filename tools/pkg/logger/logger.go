package logger

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	coreconfig "main/tools/pkg/core_config"
)

const (
	LEVEL_DEBUG   = "debug"
	LEVEL_INFO    = "info"
	LEVEL_WARNING = "warning"
	LEVEL_ERROR   = "error"
)

// Logger обертка над стандартным логгером
type Logger struct {
	*slog.Logger
}

// SetupLogger creates and configures a slog.Logger based on the specified environment.
// It takes the environment as input and returns a configured slog.Logger.
// For the 'local' environment, it uses a TextHandler, and for 'dev', it uses a JSONHandler.
func NewLogger(cfg *coreconfig.Logging) (*Logger, error) {
	var level slog.Level
	switch cfg.Level {
	case LEVEL_DEBUG:
		level = slog.LevelDebug
	case LEVEL_INFO:
		level = slog.LevelInfo
	case LEVEL_WARNING:
		level = slog.LevelWarn
	case LEVEL_ERROR:
		level = slog.LevelError
	}

	var f *os.File
	var err error
	if cfg.File != "" {
		if f, err = os.OpenFile(cfg.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return nil, err
		}
	}

	var w io.Writer
	if level == slog.LevelDebug && f != nil {
		w = io.MultiWriter(f, os.Stdout)
	} else {
		w = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(w, opts))

	return &Logger{logger}, nil
}

// With -
func (l *Logger) With(args ...any) *Logger {
	n := *l
	n.Logger = n.Logger.With(args...)
	return &n
}

// Debugf -
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Infof -
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warnf -
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Errorf -
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Panicf - use log.Panic because slog dont have this method
func (l *Logger) Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}
