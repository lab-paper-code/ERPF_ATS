package logic

import (
	"awds/types"
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
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

func (logic *Logic) GetDevice(deviceID string) (types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetDevice",
	})

	logger.Debug("received GetDevice()")

	return logic.dbAdapter.GetDevice(deviceID)
}

func (logic *Logic) GetDeviceResourceMetrics(device *types.Device) (*types.Device, error) {
	type deviceResourceMetrics struct {
		Memory			float64		`json:"memory"`
		NetworkLatency 	float64		`json:"network_latency"`
	}

	var response deviceResourceMetrics

	// get periodically stored metrics instead of measuring metrics again
	requestAddr := fmt.Sprintf("http://%s:%s/computing_measure", device.IP, device.Port)
	
	client := resty.New()
	_, err := client.R().SetResult(&response).Get(requestAddr)
	if err != nil {
		return nil, err
	}

	device.NetworkLatency = response.NetworkLatency // in Mbps
	device.Memory = response.Memory // to verify upperlimit

	return device, nil
}

func (logic *Logic) CreateDevice(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateDevice",
	})

	logger.Debug("received CreateDevice()")

	// get device info
	device, err := logic.GetDeviceResourceMetrics(device)
	if err != nil {
		return err
	}
	
	return logic.dbAdapter.InsertDevice(device)
}

func (logic *Logic) UpdateDeviceEndpoint(deviceID string, endpoint string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceIP",
	})

	logger.Debug("received UpdateDeviceIP()")

	return logic.dbAdapter.UpdateDeviceEndpoint(deviceID, endpoint)
}

func (logic *Logic) UpdateDeviceDescription(deviceID string, description string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceDescription",
	})

	logger.Debug("received UpdateDeviceDescription()")

	return logic.dbAdapter.UpdateDeviceDescription(deviceID, description)
}

func (logic *Logic) DeleteDevice(deviceID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeleteDevice",
	})

	logger.Debug("received DeleteDevice()")

	device, err := logic.GetDevice(deviceID)

	if err != nil {
		log.Error("%s does not exist, cannot delete device", device)
		return err
	}
	
	return logic.dbAdapter.DeleteDevice(deviceID)
}
