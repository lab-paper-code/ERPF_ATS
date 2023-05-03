package types

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

const (
	volumeIDPrefix string = "vol"
)

// Volume represents a volume, holding all necessary info. about volume
type Volume struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	DeviceID   string    `json:"device_id"`
	VolumeSize int64     `json:"volume_size"` // in bytes
	Mounted    bool      `json:"mounted"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// NewVolumeID creates a new Volume ID
func NewVolumeID() string {
	return fmt.Sprintf("%s_%s", volumeIDPrefix, xid.New().String())
}
