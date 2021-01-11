package configs

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config is the main struct
type Config struct {
	MainAPI *MainAPI `yaml:"config"`
}

// MainAPI holds the specified configuration for the main api
type MainAPI struct {
	Test           *TestConfig      `yaml:"test"`
	Network        *NetworkConfig   `yaml:"network"`
	DB             *DBConfig        `yaml:"db"`
	Authentication *AuthConfig      `yaml:"authentication"`
	Filestore      *FilestoreConfig `yaml:"filestore"`
}

// NetworkConfig has the http-configuration
type NetworkConfig struct {
	IP             string   `yaml:"ip"`
	Port           string   `yaml:"port"`
	WriteTimeout   int      `yaml:"write_timeout"`
	ReadTimeout    int      `yaml:"read_timeout"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

// DBConfig holds information for the DB
type DBConfig struct {
	URLs []string `yaml:"urls"`
}

// AuthConfig holds information for accessing
// the firebase authentication-service
type AuthConfig struct {
	CredentialsFile string `yaml:"credentials_file" envconfig:"AUTH_CREDENTIALS_FILE"`
	APIKey          string `yaml:"api_key" envconfig:"AUTH_API_KEY"`
}

// TestConfig holds the settings for testing the TI-API
type TestConfig struct {
	Run    bool   `yaml:"run"`
	Secret string `yaml:"secret" envconfig:"TEST_SECRET"`
}

// FilestoreConfig holds information for the filestore
type FilestoreConfig struct {
	BasePath string `yaml:"base_path"`
}

func readYAML(path string, cfg *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Missing config-file: %v", err)
	}

	if err = yaml.NewDecoder(file).Decode(cfg); err != nil {
		return fmt.Errorf("Failed to decode config-file: %v", err)
	}
	return nil
}

func readENV(cfg *Config) error {
	if err := envconfig.Process("", cfg); err != nil {
		return fmt.Errorf("Failed to process environment variables from config: %v", err)
	}
	return nil
}

// Get returns data from environment and config.yml
func Get(path string) (*Config, error) {
	var cfg Config
	if err := readYAML(path, &cfg); err != nil {
		return nil, err
	}
	if err := readENV(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
