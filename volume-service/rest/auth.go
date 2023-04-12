package rest

import (
	"github.com/gin-gonic/gin"
)

// getAdminUserAccounts returns admin user accounts
func (adapter *RESTAdapter) getAdminUserAccounts() gin.Accounts {
	users := gin.Accounts{}

	// admin
	users[adapter.config.RestAdminUsername] = adapter.config.RestAdminPassword

	return users
}
