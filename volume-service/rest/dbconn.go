package rest

import(
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/lab-paper-code/ksv/volume-service/types"
	"time"
	log "github.com/sirupsen/logrus"
)

func (service *RESTService) GetConnector() *sql.DB {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "10.106.213.189:3306",
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               "ksv",
	}
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(connector)
	return db
}

func (service *RESTService) InsertDevice(device types.Device, db *sql.DB)(int64, error) {
	var idx int64

	res, err := db.Exec("INSERT INTO device (idx, device_ip,id,pass) VALUES (NULL, ?, ?, ?)", device.IP, device.ID, device.Pwd)
	if err != nil {return idx, err}

	idx, err = res.LastInsertId()

	if err != nil {
		return idx, err
	}

	return idx, nil
}

func(service *RESTService) ReassignRows(db *sql.DB) {
	 _, err := db.Exec("SET @CNT=0")
    service.checkError(err)
	_,err = db.Exec("UPDATE device SET device.idx = @CNT:=@CNT+1;")
	service.checkError(err)

	//auto increment 초기화
	var idx string
	err = db.QueryRow("SELECT idx FROM device ORDER BY idx DESC LIMIT 1").Scan(&idx)
	service.checkError(err)

	print(idx)

	stmt, err := db.Prepare("ALTER TABLE device AUTO_INCREMENT=?")
	defer stmt.Close()

	_, err = stmt.Exec(idx)
	service.checkError(err)
	// _, err = db.Exec("ALTER TABLE device AUTO_INCREMENT=?", idx)
	// checkError(err)

}

func (service *RESTService)checkError(err error) {
    if err != nil {
        panic(err)
    }
}