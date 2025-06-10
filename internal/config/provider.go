package config

// LoadConfigFunc is a function type that loads configuration from a file
type LoadConfigFunc func(string) (*AppConfig, error)

// DefaultConfigProvider provides a default implementation for loading application configuration
type DefaultConfigProvider struct {
	fname    string
	loadFunc LoadConfigFunc
}

// NewDefaultConfigProvider creates a new DefaultConfigProvider with the default configuration file path
func NewDefaultConfigProvider() *DefaultConfigProvider {
	return NewDefaultConfigProviderWithCustomLoader("etc/config.yml", load)
}

// NewDefaultConfigProviderWithCustomLoader creates a new DefaultConfigProvider with a custom loader function
func NewDefaultConfigProviderWithCustomLoader(fname string, loader LoadConfigFunc) *DefaultConfigProvider {
	return &DefaultConfigProvider{
		fname:    fname,
		loadFunc: loader,
	}
}

// GetAppConfig provides AppConfig from default config file
func (o *DefaultConfigProvider) GetAppConfig() (res *AppConfig, err error) {
	return o.loadFunc(o.fname)
}
