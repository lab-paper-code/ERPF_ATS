package db

import (
	"awds/types"

	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListJobs() ([]types.Job, error) {
	jobs := []types.Job{}
	result := adapter.db.Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}

	return jobs, nil
}

func (adapter *DBAdapter) GetJob(jobID string) (types.Job, error) {
	var job types.Job
	result := adapter.db.Where("id = ?", jobID).First(&job)
	if result.Error != nil {
		return job, result.Error
	}

	return job, nil
}

func (adapter *DBAdapter) InsertJob(job *types.Job) error {
	result := adapter.db.Create(&job)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert an job")
	}

	return nil
}

// func (adapter *DBAdapter) UpdateAppName(appID string, name string) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.Name = name

// 	adapter.db.Save(&record)

// 	return nil
// }

// func (adapter *DBAdapter) UpdateAppRequireGPU(appID string, requireGPU bool) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.RequireGPU = requireGPU

// 	adapter.db.Save(&record)

// 	return nil
// }

// func (adapter *DBAdapter) UpdateAppDescription(appID string, description string) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.Description = description

// 	adapter.db.Save(&record)

// 	return nil
// }

// func (adapter *DBAdapter) UpdateAppDockerImage(appID string, dockerImage string) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.DockerImage = dockerImage

// 	adapter.db.Save(&record)

// 	return nil
// }

// func (adapter *DBAdapter) UpdateAppCommands(appID string, commands string) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.Commands = commands

// 	adapter.db.Save(&record)

// 	return nil
// }

// func (adapter *DBAdapter) UpdateAppArguments(appID string, arguments string) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.Arguments = arguments

// 	adapter.db.Save(&record)

// 	return nil
// }

// func (adapter *DBAdapter) UpdateAppStateful(appID string, stateful bool) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.Stateful = stateful

// 	adapter.db.Save(&record)

// 	return nil
// }

// func (adapter *DBAdapter) UpdateAppOpenPorts(appID string, openPorts string) error {
// 	var record types.AppSQLiteObj
// 	result := adapter.db.Where("id = ?", appID).Find(&record)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	record.OpenPorts = openPorts

// 	adapter.db.Save(&record)

// 	return nil
// }