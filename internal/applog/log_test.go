package applog

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meesooqa/go-tg-bnews/internal/config"
)

func TestNewConsoleLoggerProvider(t *testing.T) {
	// Test constructor with different log levels
	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}

	conf := &config.Log{}
	for _, level := range levels {
		conf.Level = level
		provider := NewConsoleLoggerProvider(conf)
		assert.Equal(t, level, provider.conf.Level, "Logger level should match the one provided")
	}
}

func TestConsoleLoggerProvider_GetLogger(t *testing.T) {
	// Redirect stdout to capture output
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	// Create logger with DEBUG level
	conf := &config.Log{Level: slog.LevelDebug}
	provider := NewConsoleLoggerProvider(conf)
	logger, cleanup := provider.GetLogger()

	// Test logging
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Restore stdout
	w.Close()
	os.Stdout = origStdout

	// Read captured output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	// Check that all messages were logged
	output := buf.String()
	assert.Contains(t, output, "debug message")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")

	// Cleanup should be a noop, but we call it anyway for completeness
	cleanup()
}

func TestConsoleLoggerProvider_GetLogger_Level(t *testing.T) {
	// Test INFO level filtering
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	// Create logger with INFO level
	conf := &config.Log{Level: slog.LevelInfo}
	provider := NewConsoleLoggerProvider(conf)
	logger, cleanup := provider.GetLogger()

	logger.Debug("debug message")
	logger.Info("info message")

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	output := buf.String()
	assert.NotContains(t, output, "debug message", "Debug messages should be filtered out at INFO level")
	assert.Contains(t, output, "info message", "Info messages should be included at INFO level")

	cleanup()
}

func TestNewFileLoggerProvider(t *testing.T) {
	// Test constructor
	level := slog.LevelDebug
	filePath := "test.log"

	conf := &config.Log{Level: level, Path: filePath}
	provider := NewFileLoggerProvider(conf)
	assert.Equal(t, level, provider.conf.Level, "Logger level should match the one provided")
	assert.Equal(t, filePath, provider.conf.Path, "Logger path should match the one provided")
}

func TestFileLoggerProvider_GetLogger(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "test.log")

	// Create logger
	conf := &config.Log{Level: slog.LevelDebug, Path: logPath}
	provider := NewFileLoggerProvider(conf)
	logger, cleanup := provider.GetLogger()

	// Write log messages
	testMessage := "file logger test"
	logger.Info(testMessage)

	// Close the file
	cleanup()

	// Read log file content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Check log content
	assert.Contains(t, string(content), testMessage)
}

func TestFileLoggerProvider_GetLogger_Level(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "test_level.log")

	// Create logger with WARN level
	conf := &config.Log{Level: slog.LevelWarn, Path: logPath}
	provider := NewFileLoggerProvider(conf)
	logger, cleanup := provider.GetLogger()

	// Write log messages of different levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Close the file
	cleanup()

	// Read log file content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Check log content - only WARN and ERROR should be present
	output := string(content)
	assert.NotContains(t, output, "debug message", "Debug messages should be filtered out at WARN level")
	assert.NotContains(t, output, "info message", "Info messages should be filtered out at WARN level")
	assert.Contains(t, output, "warn message", "Warn messages should be included at WARN level")
	assert.Contains(t, output, "error message", "Error messages should be included at WARN level")
}

func TestFileLoggerProvider_GetLogger_FileCreation(t *testing.T) {
	// Test that logger creates a new log file if it doesn't exist
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Use a subdirectory that doesn't exist yet
	nonExistentDir := filepath.Join(tempDir, "logs")
	logPath := filepath.Join(nonExistentDir, "new_file.log")

	// Create the directory
	err = os.Mkdir(nonExistentDir, 0o755)
	require.NoError(t, err)

	// Verify file doesn't exist yet
	_, err = os.Stat(logPath)
	assert.True(t, os.IsNotExist(err))

	// Create logger
	conf := &config.Log{Level: slog.LevelInfo, Path: logPath}
	provider := NewFileLoggerProvider(conf)
	logger, cleanup := provider.GetLogger()

	// Write a message
	logger.Info("new file test")

	// Close the file
	cleanup()

	// Verify file was created
	_, err = os.Stat(logPath)
	assert.NoError(t, err)

	// Check content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(content), "new file test"))
}

func TestLoggerProvider_Interface(t *testing.T) {
	// Test that both logger providers implement the LoggerProvider interface
	var providers []Provider

	conf := &config.Log{Level: slog.LevelInfo, Path: "testdata/test.log"}
	providers = append(providers, NewConsoleLoggerProvider(conf), NewFileLoggerProvider(conf))

	for i, provider := range providers {
		logger, cleanup := provider.GetLogger()
		assert.NotNil(t, logger, "Logger %d should not be nil", i)
		assert.NotNil(t, cleanup, "Cleanup function %d should not be nil", i)
		cleanup() // call cleanup to prevent resource leaks
	}
}
