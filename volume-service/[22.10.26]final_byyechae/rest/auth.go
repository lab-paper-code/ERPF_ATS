package rest

import (
	"github.com/gin-gonic/gin"
)

var (
	// id-pwd map
	authenticatedUsers = gin.Accounts{}
)

func init() {
	// this is called automatically
	authenticatedUsers["admin"] = "admin"
}

// getUserAccounts returns user accounts
func (service *RESTService) getUserAccounts() gin.Accounts {
	return authenticatedUsers
}
