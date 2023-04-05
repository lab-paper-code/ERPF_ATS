package rest

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// getUserAccounts returns user accounts
func (service *RESTService) getUserAccounts() gin.Accounts {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "getUserAccounts",
	})

	users := gin.Accounts{}

	// admin
	users[service.config.RestAdminUsername] = service.config.RestAdminPassword

	// devices
	devices, err := service.db.ListDevices()
	if err != nil {
		logger.WithError(err).Errorf("failed to list devices for authentication")
		return users
	}

	for _, device := range devices {
		users[device.ID] = device.Password
	}

	return users
}

// getAdminUserAccounts returns admin user accounts
func (service *RESTService) getAdminUserAccounts() gin.Accounts {
	users := gin.Accounts{}

	// admin
	users[service.config.RestAdminUsername] = service.config.RestAdminPassword

	return users
}
