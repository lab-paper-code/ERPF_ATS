package rest

import log "github.com/sirupsen/logrus"

// authenticateUser authenticates user
// TODO: Need to implement this
func (service *RESTService) authenticateUser(username string, password string) bool {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "authenticateUser",
	})

	logger.Debugf("authenticating a user %s", username)

	// login successful
	if username == "admin" && password == "admin" {
		logger.Debugf("authenticated a user %s", username)
		return true
	}

	logger.Debugf("failed to authenticate a user %s", username)
	return false
}
