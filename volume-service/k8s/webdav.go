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

	webdavContainerVolumeName  string = "webdav-storage"
	webdavContainerPVMountPath string = "/uploads"
)

func (adapter *K8SAdapter) GetWebdavDeploymentName(device *types.Device) string {
	return fmt.Sprintf("%s_%s", webdavDeploymentNamePrefix, device.ID)
}

func (adapter *K8SAdapter) GetWebdavDeploymentLabels(device *types.Device) map[string]string {
	labels := map[string]string{}
	labels["webdav-name"] = adapter.GetWebdavDeploymentName(device)
	labels["volume-name"] = adapter.GetVolumeName(device)
	labels["device-id"] = device.ID
	return labels
}

func (adapter *K8SAdapter) getWebdavContainers(device *types.Device) []apiv1.Container {
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

func (adapter *K8SAdapter) getWebdavContainerVolumes(device *types.Device) []apiv1.Volume {
	pvcName := adapter.GetVolumeClaimName(device)

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

func (adapter *K8SAdapter) CreateWebdavDeployment(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateWebdavDeployment",
	})

	logger.Debug("received CreateWebdavDeployment()")

	webdavDeploymentName := adapter.GetWebdavDeploymentName(device)
	webdavDeploymentLabels := adapter.GetWebdavDeploymentLabels(device)
	webdavDeploymentNumReplicas := int32(1)

	webdavContainers := adapter.getWebdavContainers(device)
	webdavContainerVolumes := adapter.getWebdavContainerVolumes(device)

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
