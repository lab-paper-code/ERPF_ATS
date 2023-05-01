package db

import (
	"github.com/lab-paper-code/ksv/volume-service/types"
	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListApps() ([]types.App, error) {
	apps := []types.App{}
	result := adapter.db.Find(&apps)
	if result.Error != nil {
		return nil, result.Error
	}

	return apps, nil
}

func (adapter *DBAdapter) GetApp(appID string) (types.App, error) {
	var app types.App
	result := adapter.db.Where("id = ?", appID).First(&app)
	if result.Error != nil {
		return app, result.Error
	}

	return app, nil
}

func (adapter *DBAdapter) InsertApp(app *types.App) error {
	result := adapter.db.Create(app)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert an app")
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
