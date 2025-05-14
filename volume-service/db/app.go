package db

import (
	"volume-service/types"

	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListApps() ([]types.App, error) {
	sqliteApps := []types.AppSQLiteObj{}
	result := adapter.db.Find(&sqliteApps)
	if result.Error != nil {
		return nil, result.Error
	}

	// convert to App
	apps := []types.App{}
	for _, sqliteApp := range sqliteApps {
		apps = append(apps, sqliteApp.ToAppObj())
	}

	return apps, nil
}

func (adapter *DBAdapter) GetApp(appID string) (types.App, error) {
	var sqliteApp types.AppSQLiteObj
	var app types.App
	result := adapter.db.Where("id = ?", appID).First(&sqliteApp)
	if result.Error != nil {
		return app, result.Error
	}

	// convert to App
	app = sqliteApp.ToAppObj()

	return app, nil
}

func (adapter *DBAdapter) InsertApp(app *types.App) error {
	// convert to AppSQLiteObj
	sqliteApp := app.ToAppSQLiteObj()

	result := adapter.db.Create(&sqliteApp)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert an app")
	}

	return nil
}

func (adapter *DBAdapter) UpdateAppName(appID string, name string) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Name = name

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppRequireGPU(appID string, requireGPU bool) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.RequireGPU = requireGPU

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppDescription(appID string, description string) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Description = description

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppDockerImage(appID string, dockerImage string) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.DockerImage = dockerImage

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppCommands(appID string, commands string) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Commands = commands

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppArguments(appID string, arguments string) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Arguments = arguments

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppStateful(appID string, stateful bool) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Stateful = stateful

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppOpenPorts(appID string, openPorts string) error {
	var record types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.OpenPorts = openPorts

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateAppRun(appRunID string, appID string, deviceID string, volumeID string) error {
	var record types.AppRun
	result := adapter.db.Where("id = ?", appRunID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.AppID = appID
	record.DeviceID = deviceID
	record.VolumeID = volumeID

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) DeleteApp(appID string) error {
	var app types.AppSQLiteObj
	result := adapter.db.Where("id = ?", appID).Delete(&app)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to delete an app")
	}

	return nil
}

func (adapter *DBAdapter) ListAppRuns(deviceID string) ([]types.AppRun, error) {
	appRuns := []types.AppRun{}
	result := adapter.db.Where("device_id = ?", deviceID).Find(&appRuns)
	if result.Error != nil {
		return nil, result.Error
	}

	return appRuns, nil
}

func (adapter *DBAdapter) ListAllAppRuns() ([]types.AppRun, error) {
	appRuns := []types.AppRun{}
	result := adapter.db.Find(&appRuns)
	if result.Error != nil {
		return nil, result.Error
	}

	return appRuns, nil
}

func (adapter *DBAdapter) GetAppRun(appRunID string) (types.AppRun, error) {
	var appRun types.AppRun
	result := adapter.db.Where("id = ?", appRunID).First(&appRun)
	if result.Error != nil {
		return appRun, result.Error
	}

	return appRun, nil
}

func (adapter *DBAdapter) InsertAppRun(appRun *types.AppRun) error {
	result := adapter.db.Create(appRun)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert an app run")
	}

	return nil
}

func (adapter *DBAdapter) UpdateAppRunTermination(appRunID string, terminated bool) error {
	var record types.AppRun
	result := adapter.db.Where("id = ?", appRunID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Terminated = terminated

	adapter.db.Save(&record)

	return nil
}
