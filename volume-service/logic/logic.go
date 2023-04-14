package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/db"
	"github.com/lab-paper-code/ksv/volume-service/k8s"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
)

type Logic struct {
	config *commons.Config

	dbAdapter  *db.DBAdapter
	k8sAdapter *k8s.K8SAdapter
}

// Start starts Logic
func Start(config *commons.Config, dbAdapter *db.DBAdapter, k8sAdapter *k8s.K8SAdapter) (*Logic, error) {
	logic := &Logic{
		config:     config,
		dbAdapter:  dbAdapter,
		k8sAdapter: k8sAdapter,
	}

	return logic, nil
}

// Stop stops Logic
func (logic *Logic) Stop() error {
	return nil
}

func (logic *Logic) ListDevices() ([]types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListDevices",
	})

	logger.Debug("received ListDevices()")

	return logic.dbAdapter.ListDevices()
}

func (logic *Logic) GetDevice(id string) (*types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetDevice",
	})

	logger.Debug("received GetDevice()")

	device, err := logic.dbAdapter.GetDevice(id)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (logic *Logic) InsertDevice(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "InsertDevice",
	})

	logger.Debug("received InsertDevice()")

	return logic.dbAdapter.InsertDevice(device)
}

func (logic *Logic) MountVolume(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "MountVolume",
	})

	logger.Debug("received MountVolume()")

	//logger.Debug("creating PV for device %s", device.ID)
	//err := logic.k8sAdapter.CreatePV(device)
	//if err != nil {
	//	return err
	//}

	logger.Debug("creating PVC for device %s", device.ID)
	err := logic.k8sAdapter.CreatePVC(device)
	if err != nil {
		return err
	}

	volumeName := logic.k8sAdapter.GetVolumeName(device)

	logger.Debug("creating Webdav Deployment for device %s, volume %s", device.ID, volumeName)
	err = logic.k8sAdapter.CreateWebdavDeployment(device)
	if err != nil {
		return err
	}

	/*
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
	return nil
}
