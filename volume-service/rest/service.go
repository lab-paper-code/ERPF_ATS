package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lab-paper-code/ksv/volume-service/commons"
	log "github.com/sirupsen/logrus"
)

type RESTService struct {
	config    *commons.Config
	router    *mux.Router
	webServer *http.Server
}

// Start starts RESTService
func Start(config *commons.Config) (*RESTService, error) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"function": "Start",
	})

	service := &RESTService{
		config:    config,
		router:    mux.NewRouter(),
		webServer: nil,
	}

	service.addHandlers()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.RestPort),
		Handler: service.router,
	}

	service.webServer = server

	fmt.Printf("Starting RESTful service at %s\n", server.Addr)
	// listen and serve in background
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error(err)
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

	fmt.Printf("Stopping the RESTful service\n")
	err := service.webServer.Close()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// basicAuth requires authenticated user
func (service *RESTService) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			pass := service.authenticateUser(username, password)
			if pass {
				// authorized user
				next.ServeHTTP(w, r)
				return
			}
		}

		// unauthorized
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
