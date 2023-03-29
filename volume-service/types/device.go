package types

import "github.com/rs/xid"

const (
	deviceIDPrefix string = "dev_"
)

// TODO: Need to separate user information out of Device
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Device represents a device, holding all necessary info. about device
type Device struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"` // TODO: need to remove this, use User struct
	Storage  string `json:"storage"`
	//Description string `json:"description,omitempty"`
}

// NewDeviceID creates a new Device ID
func NewDeviceID() string {
	return deviceIDPrefix + xid.New().String()
}
