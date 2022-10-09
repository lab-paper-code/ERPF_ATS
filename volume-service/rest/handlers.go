package rest

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// addHandlers adds http handlers
func (service *RESTService) addHandlers() {
	service.router.HandleFunc("/", service.basicAuth(service.getRootHandler)).Methods("GET")
	service.router.HandleFunc("/func1", service.basicAuth(service.getFunc1Handler)).Methods("GET")
}

func (service *RESTService) getRootHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "getRootHandler",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)

	w.Header().Set("Content-Type", "text/html")

	_, err := w.Write([]byte("called root success!"))
	if err != nil {
		logger.Error(err)
	}
}

func (service *RESTService) getFunc1Handler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "getFunc1Handler",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)

	w.Header().Set("Content-Type", "text/html")

	_, err := w.Write([]byte("called func1 success!"))
	if err != nil {
		logger.Error(err)
	}
}
