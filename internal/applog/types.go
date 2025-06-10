package applog

import "log/slog"

// Provider is an interface for providing a logger
type Provider interface {
	GetLogger() (*slog.Logger, func())
}
