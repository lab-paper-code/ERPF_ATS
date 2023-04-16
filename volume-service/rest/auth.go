package rest

import (
	"github.com/gin-gonic/gin"
)

func (adapter *RESTAdapter) isAdminUser(user string) bool {
	return adapter.config.RestAdminUsername == user
}

// getAdminUserAccounts returns admin user accounts
func (adapter *RESTAdapter) getAdminUserAccounts() gin.Accounts {
	users := gin.Accounts{}

	// admin
	users[adapter.config.RestAdminUsername] = adapter.config.RestAdminPassword

	return users
}

// getDeviceAccounts returns device accounts
func (adapter *RESTAdapter) getDeviceAccounts() gin.Accounts {
	users := gin.Accounts{}

	// admin
	users[adapter.config.RestAdminUsername] = adapter.config.RestAdminPassword

	idPasswordMap, err := adapter.logic.ListDeviceIDPasswordMap()
	if err == nil {
		for k, v := range idPasswordMap {
			users[k] = v
		}
	}

	return users
}
