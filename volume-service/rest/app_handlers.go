package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// setupAppRouter setup http request router for app
func (adapter *RESTAdapter) setupAppRouter() {
	// any devices can call these APIs
	adapter.router.GET("/apps", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleListApps)
	adapter.router.GET("/apps/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleGetApp)
	adapter.router.POST("/apps", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleCreateApp)

	adapter.router.GET("/appruns", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleListAppRuns)
	adapter.router.POST("/appruns/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleExecuteApp)
	adapter.router.GET("/appruns/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleGetAppRun)
	adapter.router.DELETE("/appruns/:id", gin.BasicAuth(adapter.getDeviceAccounts()), adapter.handleTerminateAppRun)
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
		Arguments   string `json:"arguments,omitempty"`
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
		RequireGPU:  input.RequireGPU,
		Description: input.Description,
		DockerImage: input.Description,
		Arguments:   input.Arguments,
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

func (adapter *RESTAdapter) handleExecuteApp(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleExecuteApp",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	appID := c.Param("id")

	type appExecutionRequest struct {
		// define input required
		DeviceID string `json:"device_id,omitempty"` // optional
	}

	var input appExecutionRequest

	err := c.BindJSON(&input)
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

	logger.Debugf("Executing App ID: %s", appID)

	appRun := types.AppRun{
		ID:    types.NewAppRunID(),
		AppID: appID,
	}

	if adapter.isAdminUser(user) && len(input.DeviceID) > 0 {
		// admin is trying to create a new app
		appRun.DeviceID = input.DeviceID
	} else {
		appRun.DeviceID = user
	}

	if len(appRun.DeviceID) == 0 {
		// fail
		err = xerrors.Errorf("device ID is not given")
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(appRun.AppID) == 0 {
		// fail
		err = xerrors.Errorf("app ID is not given")
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

	type appRunTerminationRequest struct {
		// define input required
		DeviceID string `json:"device_id,omitempty"` // optional
	}

	var input appRunTerminationRequest

	err := c.BindJSON(&input)
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
