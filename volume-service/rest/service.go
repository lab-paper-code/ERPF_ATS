package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/db"
	log "github.com/sirupsen/logrus"
)

type RESTService struct {
	config     *commons.Config
	address    string
	router     *gin.Engine
	httpServer *http.Server
	db         *db.DBService
}

// Start starts RESTService
func Start(config *commons.Config, dbService *db.DBService) (*RESTService, error) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"function": "Start",
	})

	addr := fmt.Sprintf(":%d", config.RestPort)
	router := gin.Default()

	service := &RESTService{
		config:  config,
		address: addr,
		router:  router,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		db: dbService,
	}

	// setup HTTP request router
	service.setupRouter()

	logger.Infof("Starting REST service at %s\n", service.address)
	// listen and serve in background
	go func() {
		err := service.httpServer.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	return service, nil
}

// Stop stops RESTService
func (service *RESTService) Stop() error {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "Stop",
	})

	logger.Infof("Stopping the REST service\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := service.httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("Stopped the REST service service\n")

	return nil
}
