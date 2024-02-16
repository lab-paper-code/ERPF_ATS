package logic

import (
	"awds/commons"
	"awds/db"
)

type Logic struct {
	config *commons.Config

	dbAdapter  *db.DBAdapter
}

// Start starts Logic
func Start(config *commons.Config, dbAdapter *db.DBAdapter) (*Logic, error) {
	logic := &Logic{
		config:     config,
		dbAdapter:  dbAdapter,
	}

	return logic, nil
}

// Stop stops Logic
func (logic *Logic) Stop() error {
	return nil
}
