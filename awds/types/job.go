package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/xerrors"
)

const (
	jobIDPrefix string = "job"
)

// Pod represents an pod, holding all necessary info. about pod
type Job struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	DeviceID string 	  `json:"device_id"`
	PodID string 		  `json:"pod_id"`

	PartitionRate float64 `json:"partition_rate"`
	DeviceStartTime time.Time `json:"device_start_time"`
	DeviceEndTime time.Time `json:"device_end_time"`
	PodStartTime time.Time `json:"pod_start_time"`
	PodEndTime time.Time   `json:"pod_end_time"`
	
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func ValidateJobID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty job id")
	}

	prefix := fmt.Sprintf("%s_", jobIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid job id - %s", id)
	}
	return nil
}

// NewAppID creates a new App ID
func NewJobID() string {
	return fmt.Sprintf("%s_%s", jobIDPrefix, xid.New().String())
}

