package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
)

type Config struct {
	Level  string
	Format string
	File   string
}

// Init initializes the global logger. Should be called once at startup.
func Init(cfg Config) error {
	var initErr error

	once.Do(func() {
		var writers []io.Writer

		// Console writer
		writers = append(writers, os.Stdout)

		// File writer
		if cfg.File != "" {
			dir := filepath.Dir(cfg.File)
			if err := os.MkdirAll(dir, 0755); err != nil {
				initErr = err
				return
			}

			file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				initErr = err
				return
			}
			writers = append(writers, file)
		}

		multiWriter := io.MultiWriter(writers...)

		// Log level
		var level slog.Level
		switch cfg.Level {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}

		opts := &slog.HandlerOptions{
			Level: level,
		}

		var handler slog.Handler
		if cfg.Format == "json" {
			handler = slog.NewJSONHandler(multiWriter, opts)
		} else {
			handler = slog.NewTextHandler(multiWriter, opts)
		}

		globalLogger = slog.New(handler)
		slog.SetDefault(globalLogger)
	})

	return initErr
}

// Get returns the global logger. Panics if Init was not called.
func Get() *slog.Logger {
	if globalLogger == nil {
		panic("logger not initialized: call logger.Init() first")
	}
	return globalLogger
}

// Debug logs at debug level
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Info logs at info level
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Warn logs at warn level
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Error logs at error level
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// Fatal logs at error level and exits
func Fatal(msg string, args ...any) {
	Get().Error(msg, args...)
	os.Exit(1)
}
