package commons

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type Config struct {
    RestPort int    `yaml:"rest_port,omitempty" json:"rest_port,omitempty" envconfig:"VOLUME_SERVICE_REST_PORT"`
    LogLevel string `yaml:"log_level,omitempty" json:"log_level,omitempty" envconfig:"VOLUME_SERVICE_LOG_LEVEL"`

    DBUsername string `yaml:"db_username,omitempty" json:"db_username,omitempty" envconfig:"VOLUME_SERVICE_DB_USERNAME"`
    DBPassword string `yaml:"db_password,omitempty" json:"db_password,omitempty" envconfig:"VOLUME_SERVICE_DB_PASSWORD"`
    DBName     string `yaml:"db_name,omitempty" json:"db_name,omitempty" envconfig:"VOLUME_SERVICE_DB_NAME"`
    DBAddress  string `yaml:"db_address,omitempty" json:"db_address,omitempty" envconfig:"VOLUME_SERVICE_DB_ADDRESS"`

	SQLitePath string `yaml:"sqlite_path,omitempty" json:"sqlite_path,omitempty" envconfig:"VOLUME_SERVICE_SQLITE_PATH"`

	CORSAllowOrigins []string `yaml:"cors_allow_origins,omitempty" json:"cors_allow_origins,omitempty"`
	CORSAllowMethods     []string `yaml:"cors_allow_methods,omitempty" json:"cors_allow_methods,omitempty"`
	CORSAllowHeaders     []string `yaml:"cors_allow_headers,omitempty" json:"cors_allow_headers,omitempty"`
	CORSAllowCredentials bool     `yaml:"cors_allow_credentials,omitempty" json:"cors_allow_credentials,omitempty"`
	CORSMaxAgeSeconds    int      `yaml:"cors_max_age_seconds,omitempty" json:"cors_max_age_seconds,omitempty"`

	InitialBatchSize int       `yaml:"initial_batch_size,omitempty"`
	PrecomputeReferenceLatencies []float64 `yaml:"precompute_reference_latencies,omitempty"`

	TemporaryOutThreshold int  `yaml:"temporary_out_threshold,omitempty"`
	MinBatchThreshold     int  `yaml:"min_batch_threshold,omitempty"`
	MaxBatchSize          int  `yaml:"max_batch_size,omitempty"`
}


func (config *Config) Validate() {
	valid := true

	if config.RestPort == 0 {
		log.Error("❌ [config] missing required field: rest_port")
		valid = false
	}
	if config.LogLevel == "" {
		log.Error("❌ [config] missing required field: log_level")
		valid = false
	}
	if config.DBUsername == "" {
		log.Error("❌ [config] missing required field: db_username")
		valid = false
	}
	if config.DBPassword == "" {
		log.Error("❌ [config] missing required field: db_password")
		valid = false
	}
	if config.DBName == "" {
		log.Error("❌ [config] missing required field: db_name")
		valid = false
	}
	if config.DBAddress == "" {
		log.Error("❌ [config] missing required field: db_address")
		valid = false
	}
	if config.SQLitePath == "" {
		log.Error("❌ [config] missing required field: sqlite_path")
		valid = false
	}
	if len(config.CORSAllowOrigins) == 0 {
		log.Warn("⚠️ [config] cors_allow_origins is empty")
	}

	if !valid {
		log.Fatal("🚫 Configuration validation failed. Please check config.yaml or environment variables.")
	} else {
		log.WithFields(log.Fields{
			"rest_port":  config.RestPort,
			"log_level":  config.LogLevel,
			"db_address": config.DBAddress,
			"db_name":    config.DBName,
			"db_user":    config.DBUsername,
			"sqlite_path": config.SQLitePath,
			"cors_allow_origins": config.CORSAllowOrigins,
		}).Info("✅ All configuration fields validated.")
	}
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
	return &Config{}
}

// NewConfigFromJSON creates Config from JSON
func newConfigFromJSON(jsonBytes []byte) (*Config, error) {
	config := GetDefaultConfig()

	err := json.Unmarshal(jsonBytes, config)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal JSON - %v", err)
	}

	config.Validate()

	return config, nil
}

// newConfigFromYAML creates Config from YAML
func newConfigFromYAML(yamlBytes []byte) (*Config, error) {
	config := GetDefaultConfig()

	err := yaml.Unmarshal(yamlBytes, config)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal YAML - %v", err)
	}

	config.Validate()

	return config, nil
}

// NewConfigFromENV creates Config from Environmental variables
func newConfigFromENV() (*Config, error) {
	config := GetDefaultConfig()

	err := envconfig.Process("", config)
	if err != nil {
		return nil, xerrors.Errorf("failed to read config from environmental variables - %v", err)
	}

	config.Validate()

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
		yjBytes, err := os.ReadFile(configFilePath)
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

		return nil, xerrors.Errorf("unreachable line")
	}

	return nil, xerrors.Errorf("unhandled configuration file - %s", configFilePath)
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
