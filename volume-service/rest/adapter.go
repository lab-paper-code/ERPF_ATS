package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/logic"
	log "github.com/sirupsen/logrus"
)

type RESTAdapter struct {
	config     *commons.Config
	address    string
	router     *gin.Engine
	httpServer *http.Server
	logic      *logic.Logic
}

// Start starts RESTAdapter
func Start(config *commons.Config, logik *logic.Logic) (*RESTAdapter, error) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"function": "Start",
	})

	addr := fmt.Sprintf(":%d", config.RestPort)
	router := gin.Default()

	adapter := &RESTAdapter{
		config:  config,
		address: addr,
		router:  router,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		logic: logik,
	}

	// setup HTTP request router
	adapter.setupRouter()

	fmt.Printf("Starting REST service at %s\n", adapter.address)
	logger.Infof("Starting REST service at %s\n", adapter.address)
	// listen and serve in background
	go func() {
		err := adapter.httpServer.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	return adapter, nil
}

// Stop stops RESTAdapter
func (adapter *RESTAdapter) Stop() error {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "Stop",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := adapter.httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
