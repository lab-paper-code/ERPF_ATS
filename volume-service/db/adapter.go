package db

import (
	"os"
	"path/filepath"

	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// we use GORM here to connect DataBase
// check out GORM examples below
// https://gorm.io/docs/create.html

const (
	SQLiteDBFileName string = "volume-service.db"
)

type DBAdapter struct {
	config *commons.Config
	db     *gorm.DB
}

func RemoveDBFile(config *commons.Config) error {
	logger := log.WithFields(log.Fields{
		"package":  "db",
		"function": "RemoveDBFile",
	})

	absPath, err := filepath.Abs(SQLiteDBFileName)
	if err != nil {
		return err
	}

	fi, err := os.Stat(absPath)
	if err == nil && !fi.IsDir() {
		// exist
		logger.Debugf("Removing db file %s", absPath)
		return os.RemoveAll(SQLiteDBFileName)
	}

	return nil
}

// Start starts DBAdapter
func Start(config *commons.Config) (*DBAdapter, error) {
	db, err := gorm.Open(sqlite.Open(SQLiteDBFileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(types.Device{})
	if err != nil {
		return nil, err
	}

	adapter := &DBAdapter{
		config: config,
		db:     db,
	}

	return adapter, nil
}

// Stop stops DBAdapter
func (adapter *DBAdapter) Stop() error {
	return nil
}

func (adapter *DBAdapter) ListDevices() ([]types.Device, error) {
	devices := []types.Device{}
	result := adapter.db.Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}

	return devices, nil
}

func (adapter *DBAdapter) GetDevice(id string) (*types.Device, error) {
	var device types.Device
	result := adapter.db.First(&device, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &device, nil
}

func (adapter *DBAdapter) InsertDevice(device *types.Device) error {
	result := adapter.db.Create([]*types.Device{device})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert a device")
	}

	return nil
}
