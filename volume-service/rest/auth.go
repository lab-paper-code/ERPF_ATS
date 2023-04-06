package rest

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// getDeviceAccounts returns device accounts
func (adapter *RESTAdapter) getDeviceAccounts() gin.Accounts {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "getDeviceAccounts",
	})

	users := gin.Accounts{}

	// devices
	devices, err := adapter.logic.ListDevices()
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
func (adapter *RESTAdapter) getAdminUserAccounts() gin.Accounts {
	users := gin.Accounts{}

	// admin
	users[adapter.config.RestAdminUsername] = adapter.config.RestAdminPassword

	return users
}
