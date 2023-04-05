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

	service.router.GET("/devices", gin.BasicAuth(service.getAdminUserAccounts()), service.handleListDevices)
	service.router.GET("/devices:id", gin.BasicAuth(service.getUserAccounts()), service.handleGetDevice)
	service.router.POST("/devices", service.handleRegisterDevice)
}

func (service *RESTService) handlePing(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "handlePing",
	})

	logger.Infof("access request to %s", c.Request.URL)

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

	logger.Infof("access request to %s", c.Request.URL)

	devices, err := service.db.ListDevices()
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type listOutput struct {
		Devices []types.Device `json:"devices"`
	}

	output := listOutput{
		Devices: devices,
	}

	// success
	c.JSON(http.StatusOK, output)
}

func (service *RESTService) handleGetDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "handleGetDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	id := c.Param("id")
	device, err := service.db.GetDevice(id)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, device)
}

func (service *RESTService) handleRegisterDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "handleRegisterDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type deviceRegistrationRequest struct {
		IP          string `json:"ip"`
		Password    string `json:"password"`
		StorageSize string `json:"storage_size"`
		Description string `json:"description,omitempty"`
	}

	var input deviceRegistrationRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storageSizeNum := types.SizeStringToNum(input.StorageSize)

	device := types.Device{
		IP:          input.IP,
		ID:          types.NewDeviceID(),
		Password:    input.Password,
		StorageSize: storageSizeNum,
		Description: input.Description,
	}

	logger.Debugf("ID: %s\tIP: %s\tPassword: %s\tStorageSize: %d", device.ID, device.IP, device.Password, storageSizeNum)

	err = service.db.InsertDevice(&device)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device.GetRedacted())
}
