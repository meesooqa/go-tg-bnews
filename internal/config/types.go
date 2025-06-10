package config

// Provider is an interface for loading application configuration
type Provider interface {
	GetAppConfig() (res *AppConfig, err error)
}
