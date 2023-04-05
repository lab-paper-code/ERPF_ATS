package commons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	defaultRestAdminUsername string = "admin"
	defaultRestAdminPassword string = "letmein"
	defaultRestPort          int    = 31200
	defaultLogLevel          string = "fatal"
	defaultDBUsername        string = "root"
	defaultDBPassword        string = "root"
	defaultDBName            string = "root"
	defaultDBAddress         string = "localhost:3306"
)

type Config struct {
	// REST related
	RestAdminUsername string `yaml:"rest_admin_username,omitempty" json:"rest_admin_username,omitempty" envconfig:"VOLUME_SERVICE_REST_ADMIN_USERNAME"`
	RestAdminPassword string `yaml:"rest_admin_password,omitempty" json:"rest_admin_password,omitempty" envconfig:"VOLUME_SERVICE_REST_ADMIN_PASSWORD"`
	RestPort          int    `yaml:"rest_port,omitempty" json:"rest_port,omitempty" envconfig:"VOLUME_SERVICE_REST_PORT"`

	// DB related
	DBUsername string `yaml:"db_username,omitempty" json:"db_username,omitempty" envconfig:"VOLUME_SERVICE_DB_USERNAME"`
	DBPassword string `yaml:"db_password,omitempty" json:"db_password,omitempty" envconfig:"VOLUME_SERVICE_DB_PASSWORD"`
	DBAddress  string `yaml:"db_address,omitempty" json:"db_address,omitempty" envconfig:"VOLUME_SERVICE_DB_ADDRESS"`
	DBName     string `yaml:"db_name,omitempty" json:"db_name,omitempty" envconfig:"VOLUME_SERVICE_DB_NAME"`

	LogLevel string `yaml:"log_level,omitempty" json:"log_level,omitempty" envconfig:"VOLUME_SERVICE_LOG_LEVEL"`
}

// GetLogLevel returns logrus log level
func (config *Config) GetLogLevel() log.Level {
	var logLevel log.Level
	err := logLevel.UnmarshalText([]byte(config.LogLevel))
	if err != nil {
		log.Errorf("failed to get log level from string %s", config.LogLevel)
		return log.InfoLevel
	}
	return logLevel
}

// GetDefaultConfig returns a default config
func GetDefaultConfig() *Config {
	return &Config{
		RestAdminUsername: defaultRestAdminUsername,
		RestAdminPassword: defaultRestAdminPassword,
		RestPort:          defaultRestPort,
		DBUsername:        defaultDBUsername,
		DBPassword:        defaultDBPassword,
		DBAddress:         defaultDBAddress,
		DBName:            defaultDBName,
		LogLevel:          defaultLogLevel,
	}
}

// NewConfigFromJSON creates Config from JSON
func newConfigFromJSON(jsonBytes []byte) (*Config, error) {
	config := GetDefaultConfig()

	err := json.Unmarshal(jsonBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON - %v", err)
	}

	return config, nil
}

// newConfigFromYAML creates Config from YAML
func newConfigFromYAML(yamlBytes []byte) (*Config, error) {
	config := GetDefaultConfig()

	err := yaml.Unmarshal(yamlBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML - %v", err)
	}

	return config, nil
}

// NewConfigFromENV creates Config from Environmental variables
func newConfigFromENV() (*Config, error) {
	config := GetDefaultConfig()

	err := envconfig.Process("", config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from environmental variables - %v", err)
	}

	return config, nil
}

// LoadConfigFile returns Config from config file path in json/yaml
func LoadConfigFile(configFilePath string) (*Config, error) {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "LoadConfigFile",
	})

	logger.Debugf("reading config file - %s", configFilePath)
	// check if it is a file or a dir
	_, err := os.Stat(configFilePath)
	if err != nil {
		return nil, err
	}

	isYaml := isYAMLFile(configFilePath)
	isJson := isJSONFile(configFilePath)

	if isYaml || isJson {
		logger.Debugf("reading YAML/JSON config file - %s", configFilePath)

		// load from YAML/JSON
		yjBytes, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return nil, err
		}

		if isYaml {
			config, err := newConfigFromYAML(yjBytes)
			if err != nil {
				return nil, err
			}
			return config, nil
		}

		if isJson {
			config, err := newConfigFromJSON(yjBytes)
			if err != nil {
				return nil, err
			}
			return config, nil
		}

		return nil, fmt.Errorf("unreachable line")
	}

	return nil, fmt.Errorf("unhandled configuration file - %s", configFilePath)
}

// LoadConfigEnv returns Config from environmental variables
func LoadConfigEnv() (*Config, error) {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "LoadConfigEnv",
	})

	logger.Debug("reading config from environment variables")

	config, err := newConfigFromENV()
	if err != nil {
		return nil, err
	}

	return config, nil
}
