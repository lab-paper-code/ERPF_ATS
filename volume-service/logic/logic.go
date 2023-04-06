package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/db"
	"github.com/lab-paper-code/ksv/volume-service/k8s"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
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

func (logic *Logic) GetDevice(id string, authkey string) (*types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetDevice",
	})

	logger.Info("received GetDevice()")

	device, err := logic.dbAdapter.GetDevice(id)
	if err != nil {
		return nil, err
	}

	if types.CheckAuthKey(device.ID, device.Password, authkey) {
		return device, nil
	}

	return nil, xerrors.Errorf("failed to get the device with id %s due to wrong password", id)
}

func (logic *Logic) InsertDevice(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "InsertDevice",
	})

	logger.Info("received InsertDevice()")

	return logic.dbAdapter.InsertDevice(device)
}
