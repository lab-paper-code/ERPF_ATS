package k8s

import (
	"context"
	"fmt"

	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	webdavDeploymentNamePrefix string = "webdav"
	webdavDeploymentNamespace  string = objectNamespace
	webdavServiceNamePrefix    string = "webdav"
	webdavServiceNamespace     string = objectNamespace

	webdavContainerVolumeName  string = "webdav-storage"
	webdavContainerPVMountPath string = "/uploads"
)

func (adapter *K8SAdapter) GetWebdavDeploymentName(volumdID string) string {
	return fmt.Sprintf("%s_%s", webdavDeploymentNamePrefix, volumdID)
}

func (adapter *K8SAdapter) GetWebdavServiceName(volumdID string) string {
	return fmt.Sprintf("%s_%s", webdavServiceNamePrefix, volumdID)
}

func (adapter *K8SAdapter) getWebdavDeploymentLabels(volume *types.Volume) map[string]string {
	labels := map[string]string{}
	labels["webdav-name"] = adapter.GetWebdavDeploymentName(volume.ID)
	labels["volume-id"] = volume.ID
	labels["device-id"] = volume.DeviceID
	return labels
}

func (adapter *K8SAdapter) getWebdavServiceLabels(volume *types.Volume) map[string]string {
	labels := map[string]string{}
	labels["webdav-name"] = adapter.GetWebdavServiceName(volume.ID)
	labels["volume-id"] = volume.ID
	labels["device-id"] = volume.DeviceID
	return labels
}

func (adapter *K8SAdapter) getWebdavContainers(device *types.Device, volume *types.Volume) []apiv1.Container {
	return []apiv1.Container{
		{
			Name:  "webdav",
			Image: "yechae/ksv-webdav:v2",
			Ports: []apiv1.ContainerPort{
				{
					ContainerPort: 80,
				},
			},
			LivenessProbe: &apiv1.Probe{
				ProbeHandler: apiv1.ProbeHandler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/",
						Port: intstr.FromInt(80),
					},
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				FailureThreshold:    3,
			},
			ReadinessProbe: &apiv1.Probe{
				ProbeHandler: apiv1.ProbeHandler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/",
						Port: intstr.FromInt(80),
					},
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				FailureThreshold:    3,
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      webdavContainerVolumeName,
					MountPath: webdavContainerPVMountPath,
				},
			},
			Env: []apiv1.EnvVar{
				{
					Name:  "BASIC_AUTH",
					Value: "True",
				},
				{
					Name:  "WEBDAV_LOGGIN",
					Value: "info",
				},
				{
					Name:  "WEBDAV_USERNAME",
					Value: device.ID, // TODO: Need to pass this through secrets
				},
				{
					Name:  "WEBDAV_PASSWORD",
					Value: device.Password, // TODO: Need to pass this through secrets
				},
			},
		},
	}
}

func (adapter *K8SAdapter) getWebdavContainerVolumes(volume *types.Volume) []apiv1.Volume {
	pvcName := adapter.GetVolumeClaimName(volume.ID)

	containerVolumes := []apiv1.Volume{
		{
			Name: webdavContainerVolumeName,
			VolumeSource: apiv1.VolumeSource{
				PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
					ReadOnly:  false,
				},
			},
		},
	}
	return containerVolumes
}

func (adapter *K8SAdapter) createWebdavDeployment(device *types.Device, volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "createWebdavDeployment",
	})

	logger.Debug("received createWebdavDeployment()")

	webdavDeploymentName := adapter.GetWebdavDeploymentName(volume.ID)
	webdavDeploymentLabels := adapter.getWebdavDeploymentLabels(volume)
	webdavDeploymentNumReplicas := int32(1)

	webdavContainers := adapter.getWebdavContainers(device, volume)
	webdavContainerVolumes := adapter.getWebdavContainerVolumes(volume)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webdavDeploymentName,
			Labels:    webdavDeploymentLabels,
			Namespace: webdavDeploymentNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &webdavDeploymentNumReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"webdav-name": webdavDeploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   webdavDeploymentName,
					Labels: webdavDeploymentLabels,
				},
				Spec: apiv1.PodSpec{
					Containers:    webdavContainers,
					Volumes:       webdavContainerVolumes,
					RestartPolicy: apiv1.RestartPolicyAlways,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(webdavDeploymentNamespace)
	_, err := deploymentclient.Get(ctx, deployment.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := deploymentclient.Create(ctx, deployment, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := deploymentclient.Update(ctx, deployment, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteWebdavDeployment(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteWebdavDeployment",
	})

	logger.Debug("received deleteWebdavDeployment()")

	webdavDeploymentName := adapter.GetWebdavDeploymentName(volumeID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(webdavDeploymentNamespace)
	err := deploymentclient.Delete(ctx, webdavDeploymentName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createWebdavService(volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "createWebdavService",
	})

	logger.Debug("received createWebdavService()")

	webdavServiceName := adapter.GetWebdavServiceName(volume.ID)
	webdavServiceLabels := adapter.getWebdavServiceLabels(volume)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webdavServiceName,
			Labels:    webdavServiceLabels,
			Namespace: webdavServiceNamespace,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:     int32(80),
					Protocol: apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": adapter.GetWebdavDeploymentName(volume.ID),
			},
			Type: apiv1.ServiceTypeClusterIP,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(webdavServiceNamespace)
	_, err := serviceClient.Get(ctx, service.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := serviceClient.Create(ctx, service, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := serviceClient.Update(ctx, service, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteWebdavService(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteWebdavService",
	})

	logger.Debug("received deleteWebdavService()")

	webdavServiceName := adapter.GetWebdavServiceName(volumeID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(webdavServiceNamespace)
	err := serviceClient.Delete(ctx, webdavServiceName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) CreateWebdav(device *types.Device, volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateWebdav",
	})

	logger.Debug("received CreateWebdav()")

	err := adapter.createWebdavDeployment(device, volume)
	if err != nil {
		return err
	}

	err = adapter.createWebdavService(volume)
	if err != nil {
		return err
	}

	/*
		//make Webdav ingress
		err = k8sClient.createWebdavIngress(input.Username, volumeID)
		if err != nil {
			panic(err)
		}
	*/
	return nil
}

func (adapter *K8SAdapter) DeleteWebdav(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "DeleteWebdav",
	})

	logger.Debug("received DeleteWebdav()")

	//TODO: Add deleteWebdavIngress

	err := adapter.deleteWebdavService(volumeID)
	if err != nil {
		return err
	}

	err = adapter.deleteWebdavDeployment(volumeID)
	if err != nil {
		return err
	}
	return nil
}
