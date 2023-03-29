package main

import (
	"fmt"
	"strconv"

	"github.com/lab-paper-code/ksv/volume-service/commons"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setCommonFlags(command *cobra.Command) {
	command.Flags().StringP("config", "c", "", "Set config file (yaml or json)")
	command.Flags().BoolP("envconfig", "e", false, "Read config from environmental variables")
	command.Flags().BoolP("version", "v", false, "Print version")
	command.Flags().BoolP("help", "h", false, "Print help")
	command.Flags().BoolP("debug", "d", false, "Enable debug mode")
}

func processFlags(command *cobra.Command) (bool, error) {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processFlags",
	})

	debug := false
	debugFlag := command.Flags().Lookup("debug")
	if debugFlag != nil {
		debugMode, err := strconv.ParseBool(debugFlag.Value.String())
		if err != nil {
			debug = false
		}

		debug = debugMode
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	helpFlag := command.Flags().Lookup("help")
	if helpFlag != nil {
		help, err := strconv.ParseBool(helpFlag.Value.String())
		if err != nil {
			help = false
		}

		if help {
			printHelp(command)
			return false, nil // stop here
		}
	}

	versionFlag := command.Flags().Lookup("version")
	if versionFlag != nil {
		version, err := strconv.ParseBool(versionFlag.Value.String())
		if err != nil {
			version = false
		}

		if version {
			printVersion(command)
			return false, nil // stop here
		}
	}

	configFlag := command.Flags().Lookup("config")
	if configFlag != nil {
		configPath := configFlag.Value.String()
		if len(configPath) > 0 {
			loadedConfig, err := commons.LoadConfigFile(configPath)
			if err != nil {
				logger.Error(err)
				return false, err // stop here
			}

			// overwrite config
			config = loadedConfig
		}
	}

	envConfigFlag := command.Flags().Lookup("envconfig")
	if envConfigFlag != nil {
		envConfig, err := strconv.ParseBool(envConfigFlag.Value.String())
		if err != nil {
			logger.Error(err)
			return false, err // stop here
		}

		if envConfig {
			loadedConfig, err := commons.LoadConfigEnv()
			if err != nil {
				logger.Error(err)
				return false, err // stop here
			}

			// overwrite config
			config = loadedConfig
		}
	}

	log.SetLevel(config.GetLogLevel())

	// prioritize command-line flag over config files
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	return true, nil // contiue
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
