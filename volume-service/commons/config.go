package commons

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type PrometheusConfig struct {
	Host  string `yaml:"host,omitempty"`
	Port  int    `yaml:"port,omitempty"`
	Query string `yaml:"query,omitempty"`
}

type VolumeConfig struct {
	DeviceCount  int   `yaml:"device_count,omitempty"`
	MinSizeBytes int64 `yaml:"min_size_bytes,omitempty"`
}

type DBConfig struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Address  string `yaml:"address,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

type Config struct {
	RestAdminUsername string `yaml:"rest_admin_username,omitempty" envconfig:"VOLUME_SERVICE_REST_ADMIN_USERNAME"`
	RestAdminPassword string `yaml:"rest_admin_password,omitempty" envconfig:"VOLUME_SERVICE_REST_ADMIN_PASSWORD"`
	RestPort          int    `yaml:"rest_port,omitempty" envconfig:"VOLUME_SERVICE_REST_PORT"`

	KubeConfigPath string `yaml:"kube_config_path,omitempty" envconfig:"VOLUME_SERVICE_KUBE_CONFIG_PATH"`
	NoKubernetes   bool   `yaml:"no_kubernetes,omitempty" envconfig:"NO_KUBERNETES"`

	LogLevel string `yaml:"log_level,omitempty" envconfig:"VOLUME_SERVICE_LOG_LEVEL"`

	DB         DBConfig         `yaml:"db,omitempty"`
	Prometheus PrometheusConfig `yaml:"prometheus,omitempty"`
	Volume     VolumeConfig     `yaml:"volume,omitempty"`
	CORSAllowedOrigins []string `yaml:"cors_allowed_origins,omitempty"`
}

var AppConfig Config

func (config *Config) GetLogLevel() log.Level {
	var logLevel log.Level
	err := logLevel.UnmarshalText([]byte(config.LogLevel))
	if err != nil {
		log.Errorf("failed to get log level from string %s", config.LogLevel)
		return log.InfoLevel
	}
	return logLevel
}

func GetDefaultConfig() *Config {
	return &Config{
		LogLevel: "info",
	}
}

func newConfigFromJSON(jsonBytes []byte) (*Config, error) {
	config := GetDefaultConfig()
	err := json.Unmarshal(jsonBytes, config)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal JSON - %v", err)
	}
	return config, nil
}

func newConfigFromYAML(yamlBytes []byte) (*Config, error) {
	config := GetDefaultConfig()
	err := yaml.Unmarshal(yamlBytes, config)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal YAML - %v", err)
	}
	return config, nil
}

func newConfigFromENV() (*Config, error) {
	config := GetDefaultConfig()
	err := envconfig.Process("", config)
	if err != nil {
		return nil, xerrors.Errorf("failed to read config from environmental variables - %v", err)
	}
	return config, nil
}

func LoadConfigFile(configFilePath string) (*Config, error) {
	logger := log.WithFields(log.Fields{
		"package":  "commons",
		"function": "LoadConfigFile",
	})

	logger.Debugf("reading config file - %s", configFilePath)
	_, err := os.Stat(configFilePath)
	if err != nil {
		return nil, err
	}

	isYaml := isYAMLFile(configFilePath)
	isJson := isJSONFile(configFilePath)

	yjBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	switch {
	case isYaml:
		return newConfigFromYAML(yjBytes)
	case isJson:
		return newConfigFromJSON(yjBytes)
	default:
		return nil, xerrors.Errorf("unhandled configuration file type for %s", configFilePath)
	}
}

func LoadConfigEnv() (*Config, error) {
	logger := log.WithFields(log.Fields{
		"package":  "commons",
		"function": "LoadConfigEnv",
	})

	logger.Debug("reading config from environment variables")
	return newConfigFromENV()
}
