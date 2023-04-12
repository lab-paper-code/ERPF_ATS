package types

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

const (
	deviceIDPrefix string = "dev"
)

// Device represents a device, holding all necessary info. about device
type Device struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	IP          string    `json:"ip"`
	Password    string    `json:"password"`
	VolumeSize  int64     `json:"volume_size"` // in bytes
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// NewDeviceID creates a new Device ID
func NewDeviceID() string {
	return fmt.Sprintf("%s_%s", deviceIDPrefix, xid.New().String())
}

func (device *Device) CheckAuthKey(authKey string) bool {
	expectedAuthKey := GetAuthKey(device.ID, device.Password)
	return expectedAuthKey == authKey
}

func (device *Device) GetRedacted() Device {
	dev := Device{}

	// copy
	dev = *device
	dev.Password = "<REDACTED>"
	return dev
}
