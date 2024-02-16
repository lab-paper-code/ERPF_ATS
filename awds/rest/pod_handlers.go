package rest

import (
	"awds/types"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// setupDeviceRouter setup http request router for device
func (adapter *RESTAdapter) setupPodRouter() {
	// any devices can call these APIs
	adapter.router.GET("/pods",adapter.handleListPods)
	adapter.router.GET("/pods/:id", adapter.handleGetPod)
	adapter.router.POST("/pods", adapter.handleRegisterPod)
	// adapter.router.PATCH("/devices/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleUpdateDevice)
	// adapter.router.DELETE("/devices/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleDeleteDevice)

	// any devices can call these APIs
}

func (adapter *RESTAdapter) handleListPods(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListPods",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type listOutput struct {
		Pods []types.Pod `json:"pods"`
	}

	output := listOutput{}

	pods, err := adapter.logic.ListPods()
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output.Pods = pods

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetPod(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetPod",
	})

	logger.Infof("access request to %s", c.Request.URL)

	podID := c.Param("id")

	err := types.ValidatePodID(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := adapter.logic.GetPod(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, device)
}

func (adapter *RESTAdapter) handleRegisterPod(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleRegisterPod",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type podRegistrationRequest struct {
		Endpoint 	string	`json:"endpoint"`
		Description string `json:"description,omitempty"`
	}

	var input podRegistrationRequest
	
	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	// if len(input.EndPoint) == 0 {
	// 	remoteAddrFields := strings.Split(c.Request.RemoteAddr, ":")
	// 	if len(remoteAddrFields) > 0 {
	// 		input.EndPoint = remoteAddrFields[0]
	// 	}
	// }

	pod := types.Pod{
		ID:          	   types.NewPodID(),
		Endpoint:          input.Endpoint,
		Description: 	   input.Description, // optional
	}

	err = adapter.logic.RegisterPod(&pod)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}

// func (adapter *RESTAdapter) handleUpdateDevice(c *gin.Context) {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "rest",
// 		"struct":   "RESTAdapter",
// 		"function": "handleUpdateDevice",
// 	})

// 	logger.Infof("access request to %s", c.Request.URL)

// 	user := c.GetString(gin.AuthUserKey)
// 	deviceID := c.Param("id")

// 	err := types.ValidateDeviceID(deviceID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if !adapter.isAdminUser(user) && deviceID != user {
// 		// requesting other's device info
// 		err := xerrors.Errorf("failed to update device %s, you cannot access other device info", deviceID)
// 		logger.Error(err)
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}

// 	type deviceUpdateRequest struct {
// 		EndPoint          string `json:"endpoint"`
// 		Description 	  string `json:"description,omitempty"`
// 	}

// 	var input deviceUpdateRequest

// 	err = c.BindJSON(&input)
// 	fmt.Println(input)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if len(input.EndPoint) > 0 {
// 		// update IP
// 		err = adapter.logic.UpdateDeviceIP(deviceID, input.EndPoint)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	if len(input.Description) > 0 {
// 		// update password
// 		err = adapter.logic.UpdateDeviceDescription(deviceID, input.Description)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	device, err := adapter.logic.GetDevice(deviceID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, device)
// }

// func (adapter *RESTAdapter) handleDeleteDevice(c *gin.Context) {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "rest",
// 		"struct":   "RESTAdapter",
// 		"function": "handleDeleteDevice",
// 	})

// 	logger.Infof("access request to %s", c.Request.URL)

// 	user := c.GetString(gin.AuthUserKey)
// 	deviceID := c.Param("id")

// 	err := types.ValidateDeviceID(deviceID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	device, err := adapter.logic.GetDevice(deviceID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if !adapter.isAdminUser(user) {
// 		err := xerrors.Errorf("failed to delete device %s, only admin can delete device", deviceID)
// 		logger.Error(err)
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}

// 	logger.Debugf("Deleting Device ID: %s", deviceID)

// 	err = adapter.logic.DeleteDevice(deviceID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, device)
// }
