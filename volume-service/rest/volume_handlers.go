package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

const (
	volumeSizeMinimum int64 = 1024 * 1024 * 1024 // 1GB
)

// setupVolumeRouter setup http request router for volume
func (adapter *RESTAdapter) setupVolumeRouter() {
	// any devices can call these APIs
	adapter.router.GET("/volumes", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleListVolumes)
	adapter.router.GET("/volumes/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleGetVolume)
	adapter.router.POST("/volumes", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleCreateVolume)
	adapter.router.PATCH("/volumes/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleUpdateVolume)

	adapter.router.POST("/mounts/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleMountVolume)
	adapter.router.DELETE("/mounts/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleUnmountVolume)
}

func (adapter *RESTAdapter) handleListVolumes(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListVolumes",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)

	type listOutput struct {
		Volumes []types.Volume `json:"volumes"`
	}

	output := listOutput{}

	if adapter.isAdminUser(user) {
		// admin - returns all volumes
		volumes, err := adapter.logic.ListAllVolumes()
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		output.Volumes = volumes
	} else {
		// device - returns mine
		volumes, err := adapter.logic.ListVolumes(user)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		output.Volumes = volumes
	}

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleCreateVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleCreateVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)

	type volumeCreationRequest struct {
		DeviceID   string `json:"device_id,omitempty"`
		VolumeSize string `json:"volume_size,omitempty"`
	}

	var input volumeCreationRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volumeSizeNum := types.SizeStringToNum(input.VolumeSize)
	if volumeSizeNum < volumeSizeMinimum {
		logger.Debugf("you cannot give volume size lesser than %d, set to %d", volumeSizeMinimum, volumeSizeMinimum)
		volumeSizeNum = volumeSizeMinimum
	}

	volume := types.Volume{
		ID:         types.NewVolumeID(),
		VolumeSize: volumeSizeNum,
	}

	if adapter.isAdminUser(user) {
		volume.DeviceID = input.DeviceID
	} else {
		volume.DeviceID = user
	}

	if len(volume.DeviceID) == 0 {
		// fail
		err = xerrors.Errorf("device ID is not given")
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("ID: %s\tVolumeSize: %d", volume.ID, volumeSizeNum)

	err = adapter.logic.CreateVolume(&volume)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleUpdateVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdateVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	type volumeUpdateRequest struct {
		VolumeSize string `json:"volume_size,omitempty"`
	}

	var input volumeUpdateRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.VolumeSize) == 0 {
		// no change
		err := xerrors.Errorf("failed to update volume %s, no change", volumeID)
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// resize
	volumeSizeNum := types.SizeStringToNum(input.VolumeSize)
	if volumeSizeNum < volumeSizeMinimum {
		logger.Debugf("you cannot give volume size lesser than %d, set to %d", volumeSizeMinimum, volumeSizeMinimum)
		volumeSizeNum = volumeSizeMinimum
	}

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("ID: %s\tVolumeSize: %d", volumeID, volumeSizeNum)

	if volume.VolumeSize == volumeSizeNum {
		// no change
		err := xerrors.Errorf("failed to resize volume %s, no size change, current %d, new %d", volumeID, volume.VolumeSize, volumeSizeNum)
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = adapter.logic.ResizeVolume(volumeID, volumeSizeNum)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleMountVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleMountVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	type volumeMountRequest struct {
		// define input required
	}

	var input volumeMountRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Mounting Volume ID: %s", volumeID)

	err = adapter.logic.MountVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleUnmountVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUnmountVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	type volumeUnmountRequest struct {
		// define input required
	}

	var input volumeUnmountRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Unmounting Volume ID: %s", volumeID)

	err = adapter.logic.UnmountVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}
