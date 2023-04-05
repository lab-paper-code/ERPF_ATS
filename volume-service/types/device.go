package types

import (
	"github.com/rs/xid"
)

const (
	deviceIDPrefix string = "dev_"
)

// Device represents a device, holding all necessary info. about device
type Device struct {
	ID          string `json:"id"`
	IP          string `json:"ip"`
	Password    string `json:"password"`
	StorageSize uint64 `json:"storage_size"`
	Description string `json:"description,omitempty"`
}

// NewDeviceID creates a new Device ID
func NewDeviceID() string {
	return deviceIDPrefix + xid.New().String()
}

func (device *Device) GetRedacted() Device {
	dev := Device{}

	// copy
	dev = *device
	dev.Password = "<REDACTED>"
	return dev
}
