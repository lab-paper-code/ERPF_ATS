package rest

import (
	"awds/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// setupAppRouter setup http request router for app
func (adapter *RESTAdapter) setupJobRouter() {
	// any devices can call these APIs
	adapter.router.GET("/jobs", adapter.handleListJobs)
	adapter.router.GET("/jobs/:id", adapter.handleGetJob)
	adapter.router.POST("/jobs", adapter.handleCreateJob)
	// adapter.router.PATCH("/jobs/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleUpdateJob)
	// adapter.router.DELETE("/jobs/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleDeleteJob)
}

func (adapter *RESTAdapter) handleListJobs(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListJobs",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type listOutput struct {
		Jobs []types.Job `json:"jobs"`
	}

	output := listOutput{}

	jobs, err := adapter.logic.ListJobs()
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output.Jobs = jobs

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetJob(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetJob",
	})

	logger.Infof("access request to %s", c.Request.URL)

	jobID := c.Param("id")

	err := types.ValidateJobID(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := adapter.logic.GetJob(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, app)
}

func (adapter *RESTAdapter) handleCreateJob(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleCreateJob",
	})

	logger.Infof("access request to %s", c.Request.URL)


	type jobCreationRequest struct {
		DeviceID string 	  `json:"device_id"`
		PodID string 		  `json:"pod_id"`

		PartitionRate float64 `json:"partition_rate"`
		DeviceStartTime time.Time `json:"device_start_time"`
		DeviceEndTime time.Time `json:"device_end_time"`
		PodStartTime time.Time `json:"pod_start_time"`
		PodEndTime time.Time   `json:"pod_end_time"`
	}

	var input jobCreationRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := types.Job{
		ID:          types.NewJobID(),
		DeviceID:        input.DeviceID,
		PodID:  		 input.PodID,
		PartitionRate: input.PartitionRate,
		DeviceStartTime: input.DeviceStartTime,
		DeviceEndTime: input.DeviceEndTime,
		PodStartTime:    input.PodStartTime,
		PodEndTime:   input.PodEndTime,
	}
	
	err = adapter.logic.CreateJob(&job)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

// func (adapter *RESTAdapter) handleUpdateJob(c *gin.Context) {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "rest",
// 		"struct":   "RESTAdapter",
// 		"function": "handleUpdateJob",
// 	})

// 	logger.Infof("access request to %s", c.Request.URL)

// 	user := c.GetString(gin.AuthUserKey)
// 	jobID := c.Param("id")

// 	err := types.ValidateJobID(jobID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	type jobUpdateRequest struct {
// 		DeviceID string 	  `json:"device_id"`
// 		PodID string 		  `json:"pod_id"`
// 		PartitionRate float64 `json:"partition_rate`
// 		DeviceStartTime time.Time `json:"device_start_time`
// 		DeviceEndTime time.Time `json:"device_end_time`
// 		PodStartTime time.Time `json:"pod_start_time`
// 		PodEndTime time.Time   `json:"pod_end_time`
// 	}

// 	var input jobUpdateRequest

// 	err = c.BindJSON(&input)
// 	fmt.Println(input)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}


// 	if input.RequireGPU {
// 		// update RequireGPU
// 		err = adapter.logic.UpdateAppRequireGPU(appID, input.RequireGPU)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	if len(input.Description) > 0 {
// 		// update Description
// 		err = adapter.logic.UpdateAppDescription(appID, input.Description)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	if len(input.DockerImage) > 0 {
// 		// update DockerImage
// 		err = adapter.logic.UpdateAppDockerImage(appID, input.DockerImage)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	if len(input.Commands) > 0 {
// 		// update Commands
// 		err = adapter.logic.UpdateAppCommands(appID, input.Commands)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	if len(input.Arguments) > 0 {
// 		// update Arguments
// 		err = adapter.logic.UpdateAppArguments(appID, input.Arguments)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	if input.Stateful {
// 		// update Stateful
// 		err = adapter.logic.UpdateAppStateful(appID, input.Stateful)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	if len(input.OpenPorts) > 0 {
// 		// update OpenPorts
// 		err = adapter.logic.UpdateAppOpenPorts(appID, input.OpenPorts)
// 		if err != nil {
// 			// fail
// 			logger.Error(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	app, err := adapter.logic.GetApp(appID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, app)
// }

// func (adapter *RESTAdapter) handleDeleteJob(c *gin.Context) {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "rest",
// 		"struct":   "RESTAdapter",
// 		"function": "handleDeleteJob",
// 	})

// 	logger.Infof("access request to %s", c.Request.URL)

// 	user := c.GetString(gin.AuthUserKey)
// 	appID := c.Param("id")

// 	err := types.ValidateAppID(appID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	app, err := adapter.logic.GetApp(appID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if !adapter.isAdminUser(user) {
// 		err := xerrors.Errorf("failed to delete app %s, only admin can delete app", appID)
// 		logger.Error(err)
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}

// 	logger.Debugf("Deleting App Run ID: %s", appID)

// 	err = adapter.logic.DeleteApp(appID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, app)
// }