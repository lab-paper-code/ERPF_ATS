package logic

import (
	"awds/types"

	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListPods() ([]types.Pod, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListPods",
	})

	logger.Debug("received ListPods()")

	return logic.dbAdapter.ListPods()
}

func (logic *Logic) GetPod(podID string) (types.Pod, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetPod",
	})

	logger.Debug("received GetPod()")

	return logic.dbAdapter.GetPod(podID)
}

func (logic *Logic) RegisterPod(pod *types.Pod) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "RegisterPod",
	})

	logger.Debug("received RegisterPod()")

	return logic.dbAdapter.InsertPod(pod)
}

// func (logic *Logic) UpdateDeviceIP(deviceID string, ip string) error {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "logic",
// 		"struct":   "Logic",
// 		"function": "UpdateDeviceIP",
// 	})

// 	logger.Debug("received UpdateDeviceIP()")

// 	return logic.dbAdapter.UpdateDeviceIP(deviceID, ip)
// }

// func (logic *Logic) UpdateDevicePassword(deviceID string, password string) error {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "logic",
// 		"struct":   "Logic",
// 		"function": "UpdateDevicePassword",
// 	})

// 	logger.Debug("received UpdateDevicePassword()")

// 	return logic.dbAdapter.UpdateDevicePassword(deviceID, password)
// }

// func (logic *Logic) UpdateDeviceDescription(deviceID string, description string) error {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "logic",
// 		"struct":   "Logic",
// 		"function": "UpdateDeviceDescription",
// 	})

// 	logger.Debug("received UpdateDeviceDescription()")

// 	return logic.dbAdapter.UpdateDeviceDescription(deviceID, description)
// }

// func (logic *Logic) DeleteDevice(deviceID string) error {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "logic",
// 		"struct":   "Logic",
// 		"function": "DeleteDevice",
// 	})

// 	logger.Debug("received DeleteDevice()")

// 	device, err := logic.GetDevice(deviceID)

// 	if err != nil {
// 		return err
// 	}

// 	return logic.dbAdapter.DeleteDevice(deviceID)
// }
