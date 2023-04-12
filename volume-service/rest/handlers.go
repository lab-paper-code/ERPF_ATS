package rest

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// setupRouter setup http request router
func (adapter *RESTAdapter) setupRouter() {
	adapter.router.GET("/ping", adapter.handlePing)

	adapter.router.GET("/devices", gin.BasicAuth(adapter.getAdminUserAccounts()), adapter.handleListDevices)
	adapter.router.GET("/devices/:id", adapter.handleGetDevice)
	adapter.router.POST("/devices", adapter.handleRegisterDevice)
}

func getAuthkeyFromRequest(c *gin.Context) string {
	authorizationKey := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(authorizationKey, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func (adapter *RESTAdapter) handlePing(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
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

func (adapter *RESTAdapter) handleListDevices(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListDevices",
	})

	logger.Infof("access request to %s", c.Request.URL)

	devices, err := adapter.logic.ListDevices()
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

func (adapter *RESTAdapter) handleGetDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	id := c.Param("id")
	authKey := getAuthkeyFromRequest(c)
	if len(authKey) > 0 {
		logger.Infof("authKey %s", authKey)
	}

	device, err := adapter.logic.GetDevice(id)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !device.CheckAuthKey(authKey) {
		err = xerrors.Errorf("failed to get device %s, wrong authorization key", id)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, device)
}

func (adapter *RESTAdapter) handleRegisterDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleRegisterDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type deviceRegistrationRequest struct {
		IP          string `json:"ip"`
		Password    string `json:"password"`
		VolumeSize  string `json:"volume_size"`
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

	volumeSizeNum := types.SizeStringToNum(input.VolumeSize)

	device := types.Device{
		IP:          input.IP,
		ID:          types.NewDeviceID(),
		Password:    input.Password,
		VolumeSize:  volumeSizeNum,
		Description: input.Description,
	}

	logger.Debugf("ID: %s\tIP: %s\tVolumeSize: %d", device.ID, device.IP, volumeSizeNum)

	err = adapter.logic.InsertDevice(&device)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device.GetRedacted())
}
