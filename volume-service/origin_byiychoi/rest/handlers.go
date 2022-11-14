package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
)

// setupRouter setup http request router
func (service *RESTService) setupRouter() {
	service.router.GET("/ping", service.handlePing)

	// require authentication
	devicesGroup := service.router.Group("/devices", gin.BasicAuth(service.getUserAccounts()))
	// /devices/
	devicesGroup.GET(".", service.handleListDevices)
}

func (service *RESTService) handlePing(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "handlePing",
	})

	logger.Infof("Page access request to %s", c.Request.URL)

	type pingOutput struct {
		Message string `json:"message"`
	}

	output := pingOutput{
		Message: "pong",
	}
	c.JSON(http.StatusOK, output)
}

func (service *RESTService) handleListDevices(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "handleListDevices",
	})

	logger.Infof("Page access request to %s", c.Request.URL)

	user := c.MustGet(gin.AuthUserKey).(string)

	type listOutput struct {
		Devices []types.Device `json:"devices"`
	}

	// dummy data
	devices := []types.Device{
		{
			ID:          types.NewDeviceID(),
			Username:    user,
			Description: "dummy description",
		},
	}

	output := listOutput{
		Devices: devices,
	}
	c.JSON(http.StatusOK, output)
}
