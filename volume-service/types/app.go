package types

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

const (
	appIDPrefix string = "app"
)

// App represents an app, holding all necessary info. about app
type App struct {
	ID          string                 `json:"id" gorm:"primaryKey"`
	RequireGPU  bool                   `json:"require_gpu,omitempty"`
	Description string                 `json:"description,omitempty"`
	DockerImage string                 `json:"docker_image"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	CreatedAt   time.Time              `json:"created_at,omitempty"`
	UpdatedAt   time.Time              `json:"updated_at,omitempty"`
}

// NewAppID creates a new App ID
func NewAppID() string {
	return fmt.Sprintf("%s_%s", appIDPrefix, xid.New().String())
}
