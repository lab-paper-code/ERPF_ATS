package db

import (
	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	SQLiteDBFileName string = "volume-service.db"
)

// we use GORM here to connect DataBase
// check out GORM examples below
// https://gorm.io/docs/create.html

type DBService struct {
	config *commons.Config
	db     *gorm.DB
}

// Start starts DBService
func Start(config *commons.Config) (*DBService, error) {
	db, err := gorm.Open(sqlite.Open(SQLiteDBFileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(types.Device{})
	if err != nil {
		return nil, err
	}

	service := &DBService{
		config: config,
		db:     db,
	}

	return service, nil
}

// Stop stops DBService
func (service *DBService) Stop() error {
	return nil
}

func (service *DBService) ListDevices() ([]types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "db",
		"struct":   "DBService",
		"function": "ListDevices",
	})

	logger.Info("received ListDevices()")

	devices := []types.Device{}
	result := service.db.Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}

	return devices, nil
}

func (service *DBService) GetDevice(id string) (*types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "db",
		"struct":   "DBService",
		"function": "GetDevice",
	})

	logger.Info("received GetDevice()")

	var device types.Device
	result := service.db.First(&device, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &device, nil
}

func (service *DBService) InsertDevice(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "db",
		"struct":   "DBService",
		"function": "InsertDevice",
	})

	logger.Info("received InsertDevice()")

	result := service.db.Create([]*types.Device{device})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert a device")
	}

	return nil
}
