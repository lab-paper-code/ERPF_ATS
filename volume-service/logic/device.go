package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListDeviceIDPasswordMap() (map[string]string, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListDeviceIDPasswordMap",
	})

	logger.Debug("received ListDeviceIDPasswordMap()")

	devices, err := logic.dbAdapter.ListDevices()
	if err != nil {
		return nil, err
	}

	idPasswordMap := map[string]string{}
	for _, device := range devices {
		idPasswordMap[device.ID] = device.Password
	}

	return idPasswordMap, nil
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

func (logic *Logic) GetDevice(deviceID string) (types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetDevice",
	})

	logger.Debug("received GetDevice()")

	return logic.dbAdapter.GetDevice(deviceID)
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

func (logic *Logic) UpdateDeviceIP(deviceID string, ip string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceIP",
	})

	logger.Debug("received UpdateDeviceIP()")

	return logic.dbAdapter.UpdateDeviceIP(deviceID, ip)
}

func (logic *Logic) UpdateDevicePassword(deviceID string, password string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDevicePassword",
	})

	logger.Debug("received UpdateDevicePassword()")

	return logic.dbAdapter.UpdateDevicePassword(deviceID, password)
}
