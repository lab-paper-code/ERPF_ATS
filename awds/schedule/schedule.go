package schedule

import (
	"awds/commons"
	"awds/db"
)

type Scheduler struct {
	config *commons.Config

	dbAdapter  *db.DBAdapter
}

// Start starts Scheduler
func Start(config *commons.Config, dbAdapter *db.DBAdapter) (*Scheduler, error) {
	scheduler := &Scheduler{
		config:     config,
		dbAdapter:  dbAdapter,
	}

	return scheduler, nil
}

// Stop stops Scheduler
func (logic *Scheduler) Stop() error {
	return nil
}
