package logic

import (
	"volume-service/types"

	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListApps() ([]types.App, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListApps",
	})

	logger.Debug("received ListApps()")

	return logic.dbAdapter.ListApps()
}

func (logic *Logic) GetApp(appID string) (types.App, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetApp",
	})

	logger.Debug("received GetApp()")

	return logic.dbAdapter.GetApp(appID)
}

func (logic *Logic) CreateApp(app *types.App) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateApp",
	})

	logger.Debug("received CreateApp()")

	return logic.dbAdapter.InsertApp(app)
}

func (logic *Logic) UpdateAppName(appID string, name string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateApp",
	})

	logger.Debug("received UpdateAppName()")

	return logic.dbAdapter.UpdateAppName(appID, name)
}

func (logic *Logic) UpdateAppRequireGPU(appID string, requireGPU bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppRequireGPU",
	})

	logger.Debug("received UpdateAppRequireGPU()")

	return logic.dbAdapter.UpdateAppRequireGPU(appID, requireGPU)
}

func (logic *Logic) UpdateAppDescription(appID string, description string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppDescription",
	})

	logger.Debug("received UpdateAppDescription()")

	return logic.dbAdapter.UpdateAppDescription(appID, description)
}

func (logic *Logic) UpdateAppDockerImage(appID string, dockerImage string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppDockerImage",
	})

	logger.Debug("received UpdateAppDockerImage()")

	return logic.dbAdapter.UpdateAppDockerImage(appID, dockerImage)
}

func (logic *Logic) UpdateAppCommands(appID string, commands string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppCommands",
	})

	logger.Debug("received UpdateAppCommands()")

	return logic.dbAdapter.UpdateAppCommands(appID, commands)
}

func (logic *Logic) UpdateAppArguments(appID string, arguments string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppArguments",
	})

	logger.Debug("received UpdateAppArguments()")

	return logic.dbAdapter.UpdateAppArguments(appID, arguments)
}

func (logic *Logic) UpdateAppStateful(appID string, stateful bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppStateful",
	})

	logger.Debug("received UpdateAppStateful()")

	return logic.dbAdapter.UpdateAppStateful(appID, stateful)
}

func (logic *Logic) UpdateAppOpenPorts(appID string, openPorts string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppOpenPorts",
	})

	logger.Debug("received UpdateAppOpenPorts()")

	return logic.dbAdapter.UpdateAppOpenPorts(appID, openPorts)
}

func (logic *Logic) DeleteApp(appID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeleteApp",
	})

	logger.Debug("received DeleteApp()")

	return logic.dbAdapter.DeleteApp(appID)
}

func (logic *Logic) ListAppRuns(deviceID string) ([]types.AppRun, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListAppRuns",
	})

	logger.Debug("received ListAppRuns()")

	return logic.dbAdapter.ListAppRuns(deviceID)
}

func (logic *Logic) ListAllAppRuns() ([]types.AppRun, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListAllAppRuns",
	})

	logger.Debug("received ListAllAppRuns()")

	return logic.dbAdapter.ListAllAppRuns()
}

func (logic *Logic) GetAppRun(appRunID string) (types.AppRun, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetAppRun",
	})

	logger.Debug("received GetAppRun()")

	return logic.dbAdapter.GetAppRun(appRunID)
}

func (logic *Logic) ExecuteApp(appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ExecuteApp",
	})

	logger.Debug("received ExecuteApp()")

	app, err := logic.GetApp(appRun.AppID)
	if err != nil {
		return err
	}

	device, err := logic.dbAdapter.GetDevice(appRun.DeviceID)
	if err != nil {
		return err
	}

	volume, err := logic.dbAdapter.GetVolume(appRun.VolumeID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.CreateApp()")
	} else {
		logger.Debugf("creating App %s for device %s, volume %s", app.Name, device.ID, volume.ID)
		err = logic.k8sAdapter.CreateApp(&device, &volume, &app, appRun)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.InsertAppRun(appRun)
}

func (logic *Logic) UpdateAppRun(appRunID string, appID string, deviceID string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateAppRun",
	})

	logger.Debug("received UpdateAppRun()")

	appRun, err := logic.dbAdapter.GetAppRun(appRunID)
	if err != nil {
		return err
	}

	app, err := logic.GetApp(appID)
	if err != nil {
		return err
	}

	device, err := logic.dbAdapter.GetDevice(deviceID)
	if err != nil {
		return err
	}

	volume, err := logic.dbAdapter.GetVolume(volumeID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.UpdateAppRun()")
	} else {
		logger.Debugf("updating AppRun %s for device %s, volume %s, app %s",
			appRun.ID, device.ID, volume.ID, app.ID)
		err = logic.k8sAdapter.UpdateAppRun(&device, &volume, &app, &appRun)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.UpdateAppRun(appRunID, appID, deviceID, volumeID)

}

func (logic *Logic) TerminateAppRun(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "TerminateAppRun",
	})

	logger.Debug("received TerminateAppRun()")

	appRun, err := logic.dbAdapter.GetAppRun(appRunID)
	if err != nil {
		return err
	}

	device, err := logic.dbAdapter.GetDevice(appRun.DeviceID)
	if err != nil {
		return err
	}

	volume, err := logic.dbAdapter.GetVolume(appRun.VolumeID)
	if err != nil {
		return err
	}

	// added to delete StatefulSet
	app, err := logic.dbAdapter.GetApp(appRun.AppID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.DeleteApp()")
	} else {
		logger.Debugf("stopping App Run %s for device %s, volume %s", appRun.ID, device.ID, volume.ID)
		err = logic.k8sAdapter.DeleteApp(appRunID, app.Stateful)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.UpdateAppRunTermination(appRunID, true)
}
