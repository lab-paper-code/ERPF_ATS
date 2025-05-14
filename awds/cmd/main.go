package main

import (
	"awds/commons"
	"awds/db"
	"awds/logic"
	"awds/rest"
	"fmt"
	"os"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config *commons.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awds [args..]",
	Short: "Distributes workload adaptively",
	Long:  `AWDS distributes workload adaptively.`,
	RunE:  processCommand,
}

func Execute() error {
	return rootCmd.Execute()
}
func processCommand(cmd *cobra.Command, args []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processCommand",
	})

	cfg, cont, err := processFlags(cmd)
	if err != nil {
		return err
	}
	if !cont {
		return nil
	}

	// 서비스 시작
	logger.Info("Starting DB Adapter...")
	dbAdapter, err := db.Start(cfg)
	if err != nil {
		return err
	}
	defer dbAdapter.Stop()

	logik, err := logic.Start(cfg, dbAdapter)
	if err != nil {
		return err
	}
	defer logik.Stop()

	logger.Info("Starting REST Adapter...")
	restAdapter, err := rest.Start(cfg, logik)
	if err != nil {
		return err
	}
	defer restAdapter.Stop()

	fmt.Println("✅ AWDS running. Press Ctrl+C to stop...")
	waitForCtrlC()
	return nil
}

func main() {
	// 로그 출력 형식 설정
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})

	// CLI 플래그 등록
	setCommonFlags(rootCmd)

	// 커맨드 실행
	if err := Execute(); err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "main",
		}).Fatal(err)
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
