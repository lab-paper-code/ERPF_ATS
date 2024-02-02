package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// setupAppRouter setup http request router for app
func (adapter *RESTAdapter) setupAppRouter() {
	// any devices can call these APIs
	adapter.router.GET("/apps", adapter.basicAuthDeviceOrAdmin, adapter.handleListApps)
	adapter.router.GET("/apps/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleGetApp)
	adapter.router.POST("/apps", adapter.basicAuthDeviceOrAdmin, adapter.handleCreateApp)
	adapter.router.PATCH("/apps/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleUpdateApp)
	adapter.router.DELETE("/apps/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleDeleteApp)

	adapter.router.GET("/appruns", adapter.basicAuthDeviceOrAdmin, adapter.handleListAppRuns)
	adapter.router.GET("/appruns/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleGetAppRun)
	adapter.router.POST("/appruns/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleExecuteApp)
	adapter.router.PATCH("/appruns/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleUpdateAppRun)
	adapter.router.DELETE("/appruns/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleTerminateAppRun)
}

func (adapter *RESTAdapter) handleListApps(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListApps",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type listOutput struct {
		Apps []types.App `json:"apps"`
	}

	output := listOutput{}

	apps, err := adapter.logic.ListApps()
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output.Apps = apps

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetApp(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetApp",
	})

	logger.Infof("access request to %s", c.Request.URL)

	appID := c.Param("id")

	err := types.ValidateAppID(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := adapter.logic.GetApp(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, app)
}

func (adapter *RESTAdapter) handleCreateApp(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleCreateApp",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)

	type appCreationRequest struct {
		Name        string `json:"name"`
		RequireGPU  bool   `json:"require_gpu,omitempty"`
		Description string `json:"description,omitempty"`
		DockerImage string `json:"docker_image"`
		Commands    string `json:"commands,omitempty"`
		Arguments   string `json:"arguments,omitempty"`
		Stateful    bool   `json:"stateful,omitempty"`
		OpenPorts   []int  `json:"open_ports,omitempty"`
	}

	var input appCreationRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) {
		// non-admin is trying to create a new app
		err := xerrors.Errorf("failed to create a new app %s, only admin can create a new app", input.Name)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	app := types.App{
		ID:          types.NewAppID(),
		Name:        input.Name,
		RequireGPU:  input.RequireGPU,
		Description: input.Description,
		DockerImage: input.DockerImage,
		Commands:    input.Commands,
		Arguments:   input.Arguments,
		Stateful:    input.Stateful,
		OpenPorts:   input.OpenPorts,
	}

	err = adapter.logic.CreateApp(&app)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

func (adapter *RESTAdapter) handleUpdateApp(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdateApp",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	appID := c.Param("id")

	err := types.ValidateAppID(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) {
		err := xerrors.Errorf("failed to update app %s, only admin can update app", appID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	type appUpdateRequest struct {
		Name        string `json:"name"`
		RequireGPU  bool   `json:"require_gpu,omitempty"`
		Description string `json:"description,omitempty"`
		DockerImage string `json:"docker_image"`
		Commands    string `json:"commands,omitempty"`
		Arguments   string `json:"arguments,omitempty"`
		Stateful    bool   `json:"stateful,omitempty"`
		OpenPorts   string `json:"open_ports,omitempty"`
	}

	var input appUpdateRequest

	err = c.BindJSON(&input)
	fmt.Println(input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Name) > 0 {
		// update Name
		err = adapter.logic.UpdateAppName(appID, input.Name)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.RequireGPU {
		// update RequireGPU
		err = adapter.logic.UpdateAppRequireGPU(appID, input.RequireGPU)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.Description) > 0 {
		// update Description
		err = adapter.logic.UpdateAppDescription(appID, input.Description)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.DockerImage) > 0 {
		// update DockerImage
		err = adapter.logic.UpdateAppDockerImage(appID, input.DockerImage)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.Commands) > 0 {
		// update Commands
		err = adapter.logic.UpdateAppCommands(appID, input.Commands)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.Arguments) > 0 {
		// update Arguments
		err = adapter.logic.UpdateAppArguments(appID, input.Arguments)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.Stateful {
		// update Stateful
		err = adapter.logic.UpdateAppStateful(appID, input.Stateful)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.OpenPorts) > 0 {
		// update OpenPorts
		err = adapter.logic.UpdateAppOpenPorts(appID, input.OpenPorts)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	app, err := adapter.logic.GetApp(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

func (adapter *RESTAdapter) handleDeleteApp(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleDeleteApp",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	appID := c.Param("id")

	err := types.ValidateAppID(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := adapter.logic.GetApp(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) {
		err := xerrors.Errorf("failed to delete app %s, only admin can delete app", appID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Deleting App Run ID: %s", appID)

	err = adapter.logic.DeleteApp(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

func (adapter *RESTAdapter) handleListAppRuns(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListAppRuns",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)

	type listOutput struct {
		AppRuns []types.AppRun `json:"app_runs"`
	}

	output := listOutput{}

	if adapter.isAdminUser(user) {
		// admin - returns all app runs
		appRuns, err := adapter.logic.ListAllAppRuns()
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		output.AppRuns = appRuns
	} else {
		err := types.ValidateDeviceID(user)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// device - returns mine
		appRuns, err := adapter.logic.ListAppRuns(user)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		output.AppRuns = appRuns
	}

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetAppRun(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetAppRun",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	appRunID := c.Param("id")

	err := types.ValidateAppRunID(appRunID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appRun, err := adapter.logic.GetAppRun(appRunID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && appRun.DeviceID != user {
		// requestiong other's app run info
		err := xerrors.Errorf("failed to get app run %s, you cannot access other devices' app run info", appRunID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, appRun)
}

func (adapter *RESTAdapter) handleUpdateAppRun(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdateAppRun",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	appRunID := c.Param("id")

	err := types.ValidateAppRunID(appRunID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type appRunUpdateRequest struct {
		DeviceID string `json:"device_id"`
		VolumeID string `json:"volume_id"`
		AppID    string `json:"app_id"`
	}

	var input appRunUpdateRequest

	err = c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = types.ValidateAppID(input.AppID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = types.ValidateVolumeID(input.VolumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && input.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get apprun %s, you cannot access other devices' volume info", appRunID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	err = types.ValidateDeviceID(input.DeviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appRun, err := adapter.logic.GetAppRun(appRunID)

	if err != nil {
		logger.Error(err)
		return
	}

	err = adapter.logic.UpdateAppRun(appRunID, appRun.AppID, appRun.DeviceID, appRun.VolumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appRun)
}

func (adapter *RESTAdapter) handleExecuteApp(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleExecuteApp",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	appID := c.Param("id")

	err := types.ValidateAppID(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type appExecuteRequest struct {
		DeviceID string `json:"device_id,omitempty"`
		VolumeID string `json:"volume_id"`
	}

	var input appExecuteRequest

	err = c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = types.ValidateVolumeID(input.VolumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = adapter.logic.GetApp(appID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	volume, err := adapter.logic.GetVolume(input.VolumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volume.ID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	appRun := types.AppRun{
		ID:       types.NewAppRunID(),
		VolumeID: input.VolumeID,
		AppID:    appID,
	}

	if adapter.isAdminUser(user) {
		appRun.DeviceID = input.DeviceID
	} else {
		appRun.DeviceID = user
	}

	err = types.ValidateDeviceID(appRun.DeviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = adapter.logic.ExecuteApp(&appRun)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appRun)
}

func (adapter *RESTAdapter) handleTerminateAppRun(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleTerminateAppRun",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	appRunID := c.Param("id")

	err := types.ValidateAppRunID(appRunID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*
		type appRunTerminationRequest struct {
		}

		var input appRunTerminationRequest

		err = c.BindJSON(&input)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	*/

	appRun, err := adapter.logic.GetAppRun(appRunID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && appRun.DeviceID != user {
		// requestiong other's app run info
		err := xerrors.Errorf("failed to get app run %s, you cannot access other devices' app run info", appRunID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Terminating App Run ID: %s", appRunID)

	err = adapter.logic.TerminateAppRun(appRunID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appRun)
}
