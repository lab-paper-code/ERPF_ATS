package types

import "github.com/rs/xid"

const (
	deviceIDPrefix string = "dev_"
)

type Device struct {
	ID          string `json:"id"`
	IP          string `json:"ip"`
	Username    string `json:"username"`
	Pwd         string `json:"pwd"`
	Storage     string `json:"storage"`
	//Description string `json:"description,omitempty"`
}

// NewDeviceID creates a new Device ID
func NewDeviceID() string {
	return deviceIDPrefix + xid.New().String()
}
