package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

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

func (logic *Logic) UpdateDeviceIP(id string, ip string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceIP",
	})

	logger.Debug("received UpdateDeviceIP()")

	// TODO: Implement this

	return xerrors.Errorf("not implemented")
}

func (logic *Logic) UpdateDevicePassword(id string, password string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDevicePassword",
	})

	logger.Debug("received UpdateDevicePassword()")

	// TODO: Implement this

	return xerrors.Errorf("not implemented")
}

func (logic *Logic) ResizeDeviceVolume(id string, size int64) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ResizeDeviceVolume",
	})

	logger.Debug("received ResizeDeviceVolume()")

	// TODO: Implement this

	return xerrors.Errorf("not implemented")
}
