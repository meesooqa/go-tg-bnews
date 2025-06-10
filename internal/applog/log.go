package applog

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/meesooqa/go-tg-bnews/internal/config"
)

// ConsoleLoggerProvider is a logger provider that writes logs to Stdout
type ConsoleLoggerProvider struct {
	conf *config.Log
}

// NewConsoleLoggerProvider creates a new ConsoleLoggerProvider with the given configuration
func NewConsoleLoggerProvider(conf *config.Log) *ConsoleLoggerProvider {
	return &ConsoleLoggerProvider{
		conf: conf,
	}
}

// GetLogger returns a logger that writes to Stdout and a cleanup function
func (o *ConsoleLoggerProvider) GetLogger() (logger *slog.Logger, cleanup func()) {
	noop := func() {
		// skip any actions
	}
	return slog.New(getLogHandler(o.conf, os.Stdout, &slog.HandlerOptions{Level: o.conf.Level})), noop
}

// FileLoggerProvider is a logger provider that writes logs to a file
type FileLoggerProvider struct {
	conf *config.Log
}

// NewFileLoggerProvider creates a new FileLoggerProvider with the given configuration
func NewFileLoggerProvider(conf *config.Log) *FileLoggerProvider {
	return &FileLoggerProvider{
		conf: conf,
	}
}

// GetLogger returns a logger that writes to a file and a cleanup function
func (o *FileLoggerProvider) GetLogger() (logger *slog.Logger, cleanup func()) {
	file, err := os.OpenFile(o.conf.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	logger = slog.New(getLogHandler(o.conf, file, &slog.HandlerOptions{Level: o.conf.Level}))
	cleanup = func() {
		file.Close() // nolint
	}
	return logger, cleanup
}

func getLogHandler(conf *config.Log, w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	switch conf.OutputFormat {
	case "json":
		return slog.NewJSONHandler(w, opts)
	default:
		return slog.NewTextHandler(w, opts)
	}
}
