package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/db"
	"github.com/lab-paper-code/ksv/volume-service/k8s"
	"github.com/lab-paper-code/ksv/volume-service/rest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config *commons.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "volume-service [args..]",
	Short: "Provides controls of volumes via REST interface",
	Long:  `Volume-Serivce provides controls of volumes via REST interface.`,
	RunE:  processCommand,
}

func Execute() error {
	return rootCmd.Execute()
}

func processCommand(command *cobra.Command, args []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processCommand",
	})

	cont, err := processFlags(command)
	if err != nil {
		logger.Error(err)
	}

	if !cont {
		return err
	}

	// start service
	logger.Info("Starting DB Service...")
	dbService, err := db.Start(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer dbService.Stop()
	logger.Info("DB Service Started")

	logger.Info("Starting K8S Service...")
	k8sService, err := k8s.Start(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer k8sService.Stop()
	logger.Info("K8S Service Started")

	logger.Info("Starting REST Service...")
	restService, err := rest.Start(config, dbService, k8sService)
	if err != nil {
		logger.Fatal(err)
	}
	defer restService.Stop()
	logger.Info("REST Service Started")

	// wait
	fmt.Println("Press Ctrl+C to stop...")
	waitForCtrlC()

	return nil
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})

	config = commons.GetDefaultConfig()
	log.SetLevel(config.GetLogLevel())

	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "main",
	})

	// attach common flags
	setCommonFlags(rootCmd)

	err := Execute()
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup

	endWaiter.Add(1)
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	endWaiter.Wait()
}
