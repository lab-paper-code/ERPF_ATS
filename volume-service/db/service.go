package db

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
)

type DBService struct {
	config *commons.Config
}

// Start starts DBService
func Start(config *commons.Config) (*DBService, error) {
	service := &DBService{
		config: config,
	}

	// do connection pooling if you need

	return service, nil
}

// Stop stops DBService
func (service *DBService) Stop() error {
	return nil
}

func (service *DBService) getConnector() *sql.DB {
	cfg := mysql.Config{
		User:                 service.config.DBUsername,
		Passwd:               service.config.DBPassword,
		Net:                  "tcp",
		Addr:                 service.config.DBAddress,
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               service.config.DBName,
	}
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(connector)
	return db
}

func (service *DBService) ListDevices() ([]types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "db",
		"struct":   "DBService",
		"function": "ListDevices",
	})

	logger.Info("received ListDevices()")

	// TODO: implement this
	return nil, nil
}

func (service *DBService) GetDevice(id string) (*types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "db",
		"struct":   "DBService",
		"function": "GetDevice",
	})

	logger.Info("received GetDevice()")

	// TODO: implement this
	device := types.Device{
		ID: id,
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

	// TODO: implement this
	return nil
}

/*
func (service *DBService) InsertDevice(device types.Device, db *sql.DB) (int64, error) {
	var idx int64

	res, err := db.Exec("INSERT INTO device (idx, device_ip,id,pass) VALUES (NULL, ?, ?, ?)", device.IP, device.ID, device.Password)
	if err != nil {
		return idx, err
	}

	idx, err = res.LastInsertId()

	if err != nil {
		return idx, err
	}

	return idx, nil
}

func (service *DBService) ReassignRows(db *sql.DB) {
	_, err := db.Exec("SET @CNT=0")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("UPDATE device SET device.idx = @CNT:=@CNT+1;")
	if err != nil {
		panic(err)
	}

	//auto increment 초기화
	var idx string
	err = db.QueryRow("SELECT idx FROM device ORDER BY idx DESC LIMIT 1").Scan(&idx)
	if err != nil {
		panic(err)
	}

	print(idx)

	stmt, err := db.Prepare("ALTER TABLE device AUTO_INCREMENT=?")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(idx)
	if err != nil {
		panic(err)
	}
	// _, err = db.Exec("ALTER TABLE device AUTO_INCREMENT=?", idx)
	// checkError(err)

}
*/
