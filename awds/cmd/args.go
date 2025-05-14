package main

import (
	"awds/commons"
	"awds/db"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configPath   string
	envConfig    bool
	version      bool
	help         bool
	debug        bool
	clearDB      bool
)

func setCommonFlags(command *cobra.Command) {
	command.Flags().StringVarP(&configPath, "config", "c", "", "Set config file (yaml or json)")
	command.Flags().BoolVarP(&envConfig, "envconfig", "e", false, "Read config from environmental variables")
	command.Flags().BoolVarP(&version, "version", "v", false, "Print version")
	command.Flags().BoolVarP(&help, "help", "h", false, "Print help")
	command.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	command.Flags().BoolVar(&clearDB, "clear_db", false, "Clear DB data")
}

func processFlags(command *cobra.Command) (*commons.Config, bool, error) {
	logger := log.WithFields(log.Fields{
		"package":  "cmd",
		"function": "processFlags",
	})

	if help {
		_ = printHelp(command)
		return nil, false, nil
	}

	if version {
		_ = printVersion(command)
		return nil, false, nil
	}

	var cfg *commons.Config
	var err error

	switch {
	case configPath != "":
		cfg, err = commons.LoadConfigFile(configPath)
	case envConfig:
		cfg, err = commons.LoadConfigEnv()
	default:
		cfg, err = commons.LoadConfigFile("config.yaml")
	}
	if err != nil {
		logger.Error(err)
		return nil, false, err
	}

	// 로그 레벨 설정
	log.SetLevel(cfg.GetLogLevel())
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if clearDB {
		if err := db.RemoveDBFile(cfg); err != nil {
			logger.Error(err)
			return nil, false, err
		}
	}

	return cfg, true, nil
}

func printVersion(command *cobra.Command) error {
	info, err := commons.GetVersionJSON()
	if err != nil {
		return err
	}

	fmt.Println(info)
	return nil
}

func printHelp(command *cobra.Command) error {
	return command.Usage()
}
