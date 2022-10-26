package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"github.com/lab-paper-code/ksv/volume-service/k8s"
	"fmt"
	//"database/sql"
	//"github.com/go-sql-driver/mysql"
	//"time"
)

// setupRouter setup http request router
func (service *RESTService) setupRouter() {
	service.router.GET("/ping", service.handlePing)

	// require authentication
	devicesGroup := service.router.Group("/devices", gin.BasicAuth(service.getUserAccounts()))
	// /devices/
	devicesGroup.GET(".", service.handleListDevices)
	//service.router.GET("/register", service.handleRegister)
	service.router.POST("/register", service.handleRegister)
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
			Pwd: 		 "test",

		},
	}

	output := listOutput{
		Devices: devices,
	}
	c.JSON(http.StatusOK, output)
}

func (service *RESTService) handleRegister(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTService",
		"function": "handleRegister(device)",
	})

	logger.Infof("Page access request to %s", c.Request.URL)

	var input types.Device

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	// query_storage := c.Query("storage")  
	//TODO: 현재는 pvc 용량이 default 20GB 으로 되어 있지만,
    // 향후 storage도 parameter로 받아서 현재 서버에 생성할 수 있는 용량인지 판단할 수 있도록 

	device := types.Device{
		
		IP:          input.IP,
		ID:          "",
		Username:    input.Username,
		Pwd: 		 input.Pwd,
		Storage:     input.Storage,
}


	fmt.Println("%s    %s    %s", input.IP, input.Username, input.Pwd)
	fmt.Println(input.Storage)
	

	db := service.GetConnector()
	err := db.Ping()
	if err != nil {
		panic(err)
	}

	//TODO: 중복체크
	idx, err := service.InsertDevice(device, db)
	if err != nil {
		panic(err)
	}
	print(idx)

	volumeID := fmt.Sprintf("%s%d", "ksv", idx)

	// //pvc
	
	k8sClient, err := k8s.NewK8sClient("/home/palisade2/.kube/config")
	service.checkError(err)

	//make PVC
	err = k8sClient.CreatePVC(input.Username, volumeID)
	service.checkError(err)

	//make webdav deploy
	err = k8sClient.CreateWebdavDeploy(input.Username, volumeID)
	service.checkError(err)

	//make App deploy
	err = k8sClient.CreateAppDeploy(input.Username, volumeID)
	service.checkError(err)

	//make webdav service
	err = k8sClient.CreateWebdavSVC(input.Username, volumeID)
	service.checkError(err)

	//make App service 
	err = k8sClient.CreateAppSVC(input.Username, volumeID)
	service.checkError(err)

	//make Webdav ingress
	err = k8sClient.CreateWebdavIngress(input.Username, volumeID)
	service.checkError(err)

	//make App ingress
	err = k8sClient.CreateAppIngress(input.Username, volumeID)
	service.checkError(err)

	err = k8sClient.WaitPodRun3(input.Username, volumeID)
	service.checkError(err)
	fmt.Printf("\nAll pods in podname=\"%s\" are running!\n", volumeID)

	execCommand := "sed -i -e 's#Alias /uploads \"/uploads\"#Alias /"+volumeID+"/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf"
	//change webdav path using volumeID
	err = k8sClient.ExecInPod("vd", volumeID,execCommand)
	service.checkError(err)

	execCommand = "/usr/sbin/httpd -k restart"
	err = k8sClient.ExecInPod("vd", volumeID,execCommand)
	service.checkError(err)
	
	//TODO:  k8s resource들 생성한 후 
	//1. webdav pod으로 exec 명령어로 sed -i -e 's#Alias /uploads \"/uploads\"#Alias /<volumdID>/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf 명령어 실행
	//2. app pod으로 http://ip:60000/hello_flask?ip=<ip> 해서 dom ip 알려주기 

	type Output struct {
		Mount string `json:mountPath`
		Device types.Device `json: device`
	}

	Mount := "http://155.230.36.27/"+volumeID+"/uploads"

	device =types.Device{
			IP:          input.IP,
			ID:          volumeID,
			Username:    input.Username,
			Pwd: 		 input.Pwd,
			Storage:     input.Storage,
	}

	output := Output{
		Mount: Mount,
		Device: device,
	}

	c.JSON(http.StatusOK, output)
}



// //TODO: handleGetVolumeIDByIP 완성
// func (service *RESTService) handleGetVolumeIDByIP(c *gin.Context){

// }

// //TODO: handleGetIPByVolumeID 완성
// func (service *RESTService) handleGetIPByVolumeID(c *gin.Context){

// }

