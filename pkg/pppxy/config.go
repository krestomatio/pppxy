package pppxy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// PPPxyConfig holds the configuration for a single proxy protocol proxy instance.
type PPPxyConfig struct {
	ListenAddr           string `yaml:"listen_addr"`
	BackendAddr          string `yaml:"backend_addr"`
	ProxyProtocolVersion int    `yaml:"proxy_protocol_version"`
}

// Config holds the entire application configuration.
type Config struct {
	PPPxyGroup []PPPxyConfig `yaml:"pppxy_group"`
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
