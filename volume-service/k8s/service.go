package k8s

import (
	"github.com/lab-paper-code/ksv/volume-service/commons"
	log "github.com/sirupsen/logrus"
)

type K8SService struct {
	config *commons.Config
}

// Start starts K8SService
func Start(config *commons.Config) (*K8SService, error) {
	service := &K8SService{
		config: config,
	}

	return service, nil
}

// Stop stops K8SService
func (service *K8SService) Stop() error {
	return nil
}

func (service *K8SService) CreatePV(volID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SService",
		"function": "CreatePV",
	})

	logger.Info("received CreatePV()")

	/*
		volumeID := fmt.Sprintf("%s%d", "ksv", idx)

		// //pvc

		k8sClient, err := k8s.NewK8sClient("/home/palisade2/.kube/config")
		if err != nil {
			panic(err)
		}

		//make PVC
		err = k8sClient.CreatePVC(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		//make webdav deploy
		err = k8sClient.CreateWebdavDeploy(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		//make App deploy
		err = k8sClient.CreateAppDeploy(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		//make webdav service
		err = k8sClient.CreateWebdavSVC(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		//make App service
		err = k8sClient.CreateAppSVC(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		//make Webdav ingress
		err = k8sClient.CreateWebdavIngress(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		//make App ingress
		err = k8sClient.CreateAppIngress(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		err = k8sClient.WaitPodRun3(input.Username, volumeID)
		if err != nil {
			panic(err)
		}

		logger.Infof("All pods in podname=\"%s\" are running!", volumeID)

		execCommand := "sed -i -e 's#Alias /uploads \"/uploads\"#Alias /" + volumeID + "/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf"
		//change webdav path using volumeID
		err = k8sClient.ExecInPod("vd", volumeID, execCommand)
		if err != nil {
			panic(err)
		}

		execCommand = "/usr/sbin/httpd -k restart"
		err = k8sClient.ExecInPod("vd", volumeID, execCommand)
		if err != nil {
			panic(err)
		}

		//TODO:  k8s resource들 생성한 후
		//1. webdav pod으로 exec 명령어로 sed -i -e 's#Alias /uploads \"/uploads\"#Alias /<volumdID>/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf 명령어 실행
		//2. app pod으로 http://ip:60000/hello_flask?ip=<ip> 해서 dom ip 알려주기

		type Output struct {
			Mount  string       `json:mountPath`
			Device types.Device `json: device`
		}

		Mount := "http://155.230.36.27/" + volumeID + "/uploads"

		device = types.Device{
			IP:       input.IP,
			ID:       volumeID,
			Username: input.Username,
			Password: input.Password,
			Storage:  input.Storage,
		}

		output := Output{
			Mount:  Mount,
			Device: device,
		}
	*/

	// TODO: implement this
	return nil
}
