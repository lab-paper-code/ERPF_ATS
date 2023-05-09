package types

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

const (
	appIDPrefix string = "app"
	runIDPrefix string = "run"
)

// App represents an app, holding all necessary info. about app
type App struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	RequireGPU  bool      `json:"require_gpu,omitempty"`
	Description string    `json:"description,omitempty"`
	DockerImage string    `json:"docker_image"`
	Arguments   string    `json:"arguments,omitempty"` // a space-separated command-line arguments to run app, array/map not supported
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// NewAppID creates a new App ID
func NewAppID() string {
	return fmt.Sprintf("%s_%s", appIDPrefix, xid.New().String())
}

// AppRun represents an app execution (run), holding all necessary info. about app execution
type AppRun struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	DeviceID   string    `json:"device_id"`
	AppID      string    `json:"app_id"`
	Terminated bool      `json:"terminated"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// NewAppRunID creates a new AppRun ID
func NewAppRunID() string {
	return fmt.Sprintf("%s_%s", runIDPrefix, xid.New().String())
}
