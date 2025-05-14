package k8s

import (
	"context"
	"fmt"
	"strings"

	"volume-service/types"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

const (
	appDeploymentNamePrefix  string = "app"
	appDeploymentNamespace   string = objectNamespace
	appStatefulSetNamePrefix string = "app"
	appStatefulSetNamespace  string = objectNamespace
	appServiceNamePrefix     string = "app"
	appServiceNamespace      string = objectNamespace
	appIngressNamePrefix     string = "app"
	appIngressNamespace      string = objectNamespace

	appContainerVolumeName  string = "app-storage"
	appContainerPVMountPath string = "/uploads"

	truePointer				bool = true
)

func (adapter *K8SAdapter) GetAppDeploymentName(appRunID string) string {
	return makeValidObjectName(appDeploymentNamePrefix, appRunID)
}

func (adapter *K8SAdapter) GetAppStatefulSetName(appRunID string) string {
	return makeValidObjectName(appDeploymentNamePrefix, appRunID)
}

func (adapter *K8SAdapter) GetAppServiceName(appRunID string) string {
	return makeValidObjectName(appServiceNamePrefix, appRunID)
}

func (adapter *K8SAdapter) GetAppIngressName(appRunID string) string {
	return makeValidObjectName(appIngressNamePrefix, appRunID)
}

func (adapter *K8SAdapter) getAppDeploymentLabels(appRun *types.AppRun) map[string]string {
	labels := map[string]string{}
	labels["app-name"] = adapter.GetAppDeploymentName(appRun.ID)
	labels["app-id"] = appRun.AppID
	labels["apprun-id"] = appRun.ID
	labels["volume-id"] = appRun.VolumeID
	labels["device-id"] = appRun.DeviceID
	return labels
}

func (adapter *K8SAdapter) getAppStatefulSetLabels(appRun *types.AppRun) map[string]string {
	labels := map[string]string{}
	labels["app-name"] = adapter.GetAppStatefulSetName(appRun.ID)
	labels["app-id"] = appRun.AppID
	labels["apprun-id"] = appRun.ID
	labels["volume-id"] = appRun.VolumeID
	labels["device-id"] = appRun.DeviceID
	return labels
}

func (adapter *K8SAdapter) getAppServiceLabels(appRun *types.AppRun) map[string]string {
	labels := map[string]string{}
	labels["app-name"] = adapter.GetAppServiceName(appRun.ID)
	labels["app-id"] = appRun.AppID
	labels["apprun-id"] = appRun.ID
	labels["volume-id"] = appRun.VolumeID
	labels["device-id"] = appRun.DeviceID
	return labels
}

func (adapter *K8SAdapter) getAppIngressLabels(appRun *types.AppRun) map[string]string {
	labels := map[string]string{}
	labels["app-name"] = adapter.GetAppIngressName(appRun.ID)
	labels["app-id"] = appRun.AppID
	labels["apprun-id"] = appRun.ID
	labels["volume-id"] = appRun.VolumeID
	labels["device-id"] = appRun.DeviceID
	return labels
}

func (adapter *K8SAdapter) GetAppIngressPath(appRunID string) string {
	return fmt.Sprintf("/app/%s", appRunID)
}

func (adapter *K8SAdapter) getAppContainers(app *types.App, device *types.Device, volume *types.Volume) []apiv1.Container {
	containerPorts := []apiv1.ContainerPort{}
	for _, port := range app.OpenPorts {
		containerPorts = append(containerPorts, apiv1.ContainerPort{
			Name:          fmt.Sprintf("cont-port-%d", port),
			ContainerPort: int32(port),
		})
	}

	cmdString := app.Commands
	commands := strings.Split(cmdString, " ")

	argString := app.Arguments
	arguments := strings.Split(argString, " ")

	gpuFlag := "0"
	// set to 1 if app requires GPU
	if app.RequireGPU {
		gpuFlag = "1"
	}

	// Create a container object
	container := apiv1.Container{
		Name:            "app",
		Image:           app.DockerImage,
		ImagePullPolicy: "IfNotPresent",
		Ports:           containerPorts,
		VolumeMounts: []apiv1.VolumeMount{
			{
				Name:      appContainerVolumeName,
				MountPath: appContainerPVMountPath,
			},
		},
		SecurityContext: &apiv1.SecurityContext{
			Privileged: ptr.To[bool](true),
		},
	}

	// Conditionally set Command and Args if they are not empty
	if cmdString != "" {
		container.Command = commands
	}

	if argString != "" {
		container.Args = arguments
	}

	// Conditionally set GPU Limits if gpuFlag is not "0"
	if gpuFlag != "0" {
		container.Resources = apiv1.ResourceRequirements{
			Limits: apiv1.ResourceList{
				"nvidia.com/gpu": resource.MustParse(gpuFlag),
			},
		}
	}
	
	return []apiv1.Container{container}
}

func (adapter *K8SAdapter) getAppContainerVolumes(volume *types.Volume) []apiv1.Volume {
	pvcName := adapter.GetVolumeClaimName(volume.ID)

	containerVolumes := []apiv1.Volume{
		{
			Name: appContainerVolumeName,
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

func (adapter *K8SAdapter) createAppDeployment(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "createAppDeployment",
	})

	logger.Debug("received createAppDeployment()")

	appDeploymentName := adapter.GetAppDeploymentName(appRun.ID)
	appDeploymentLabels := adapter.getAppDeploymentLabels(appRun)
	deployReplicas := int32(1)

	appContainers := adapter.getAppContainers(app, device, volume)
	appContainerVolumes := adapter.getAppContainerVolumes(volume)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appDeploymentName,
			Labels:    appDeploymentLabels,
			Namespace: appDeploymentNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app-name": appDeploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   appDeploymentName,
					Labels: appDeploymentLabels,
				},
				Spec: apiv1.PodSpec{
					Containers:    appContainers,
					Volumes:       appContainerVolumes,
					RestartPolicy: apiv1.RestartPolicyAlways,
				}, //spec
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(appDeploymentNamespace)
	_, err := deploymentclient.Get(ctx, deployment.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist -> create
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

func (adapter *K8SAdapter) updateAppDeployment(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "updateAppDeployment",
	})

	logger.Debug("received updateAppDeployment()")

	appDeploymentName := adapter.GetAppDeploymentName(appRun.ID)
	appDeploymentLabels := adapter.getAppDeploymentLabels(appRun)
	// update labels
	appDeploymentLabels["app-id"] = app.ID
	appDeploymentLabels["device-id"] = device.ID
	appDeploymentLabels["volume-id"] = volume.ID

	appContainers := adapter.getAppContainers(app, device, volume)
	appContainerVolumes := adapter.getAppContainerVolumes(volume)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(appDeploymentNamespace)
	existingDeployment, err := deploymentclient.Get(ctx, appDeploymentName, metav1.GetOptions{})

	existingDeployment.ObjectMeta.Labels = appDeploymentLabels
	existingDeployment.Spec.Template.ObjectMeta.Labels = appDeploymentLabels
	existingDeployment.Spec.Template.Spec.Containers = appContainers
	existingDeployment.Spec.Template.Spec.Volumes = appContainerVolumes

	// does not exist -> cannot update
	if err != nil {
		return err
	}

	// exist -> update
	_, updateErr := deploymentclient.Update(ctx, existingDeployment, metav1.UpdateOptions{})
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (adapter *K8SAdapter) deleteAppDeployment(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteAppDeployment",
	})

	logger.Debug("received deleteAppDeployment()")

	appDeploymentName := adapter.GetAppDeploymentName(appRunID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(appDeploymentNamespace)
	err := deploymentclient.Delete(ctx, appDeploymentName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createAppStatefulSet(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "createAppStatefulSet",
	})

	logger.Debug("received createAppStatefulSet()")

	appStatefulSetName := adapter.GetAppStatefulSetName(appRun.ID)
	appStatefulSetLabels := adapter.getAppStatefulSetLabels(appRun)
	stsReplicas := int32(1)

	appContainers := adapter.getAppContainers(app, device, volume)
	appContainerVolumes := adapter.getAppContainerVolumes(volume)

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appStatefulSetName,
			Labels:    appStatefulSetLabels,
			Namespace: appStatefulSetNamespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &stsReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app-name": appStatefulSetName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   appStatefulSetName,
					Labels: appStatefulSetLabels,
				},
				Spec: apiv1.PodSpec{
					Containers:    appContainers,
					Volumes:       appContainerVolumes,
					RestartPolicy: apiv1.RestartPolicyAlways,
				}, //spec
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	statefulsetclient := adapter.clientSet.AppsV1().StatefulSets(appStatefulSetNamespace)
	_, err := statefulsetclient.Get(ctx, statefulset.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist -> create
		_, createErr := statefulsetclient.Create(ctx, statefulset, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else{
	// exist -> update
	_, updateErr := statefulsetclient.Update(ctx, statefulset, metav1.UpdateOptions{})
	if updateErr != nil {
		return updateErr
	}
	}

	return nil
}

func (adapter *K8SAdapter) updateAppStatefulSet(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "updateAppStatefulSet",
	})

	logger.Debug("received updateAppStatefulSet()")

	appStatefulSetName := adapter.GetAppStatefulSetName(appRun.ID)
	appStatefulSetLabels := adapter.getAppStatefulSetLabels(appRun)
	// update labels
	appStatefulSetLabels["app-id"] = app.ID
	appStatefulSetLabels["device-id"] = device.ID
	appStatefulSetLabels["volume-id"] = volume.ID

	stsReplicas := int32(1)

	appContainers := adapter.getAppContainers(app, device, volume)
	appContainerVolumes := adapter.getAppContainerVolumes(volume)

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appStatefulSetName,
			Labels:    appStatefulSetLabels,
			Namespace: appStatefulSetNamespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &stsReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app-name": appStatefulSetName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   appStatefulSetName,
					Labels: appStatefulSetLabels,
				},
				Spec: apiv1.PodSpec{
					Containers:    appContainers,
					Volumes:       appContainerVolumes,
					RestartPolicy: apiv1.RestartPolicyAlways,
				}, //spec
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	statefulsetclient := adapter.clientSet.AppsV1().StatefulSets(appStatefulSetNamespace)
	_, err := statefulsetclient.Get(ctx, statefulset.GetName(), metav1.GetOptions{})
	// does not exist -> cannot update
	if err != nil {
		return err
	} else {
		// exist -> update
		_, updateErr := statefulsetclient.Update(ctx, statefulset, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteAppStatefulSet(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteAppStatefulSet",
	})

	logger.Debug("received deleteAppStatefulSet()")

	appStatefulSetName := adapter.GetAppStatefulSetName(appRunID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	statefulsetclient := adapter.clientSet.AppsV1().StatefulSets(appStatefulSetNamespace)
	err := statefulsetclient.Delete(ctx, appStatefulSetName, *metav1.NewDeleteOptions(0)) // TODO: check how to gracefully delete staefulset
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createAppService(app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "createAppService",
	})

	logger.Debug("received createAppService()")

	appServiceName := adapter.GetAppServiceName(appRun.ID)
	appServiceLabels := adapter.getAppServiceLabels(appRun)

	servicePorts := []apiv1.ServicePort{}
	for _, port := range app.OpenPorts {
		servicePorts = append(servicePorts, apiv1.ServicePort{
			Name:     fmt.Sprintf("svc-port-%d", port),
			Port:     int32(port),
			Protocol: apiv1.ProtocolTCP,
		})
	}

	// prepare selector based on stateful or stateless app
	selector := map[string]string{
		"app-name": adapter.GetAppDeploymentName(appRun.ID),
	}
	if app.Stateful {
		selector = map[string]string{
			"app-name": adapter.GetAppStatefulSetName(appRun.ID),
		}
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appServiceName,
			Labels:    appServiceLabels,
			Namespace: appServiceNamespace,
		},
		Spec: apiv1.ServiceSpec{
			Ports:    servicePorts,
			Selector: selector,
			Type:     apiv1.ServiceTypeClusterIP,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(volumeNamespace)
	_, err := serviceClient.Get(ctx, service.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist -> create
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

func (adapter *K8SAdapter) updateAppService(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "updateAppService",
	})

	logger.Debug("received updateAppService()")

	appServiceName := adapter.GetAppServiceName(appRun.ID)
	appServiceLabels := adapter.getAppServiceLabels(appRun)
	// update labels
	appServiceLabels["app-id"] = app.ID
	appServiceLabels["device-id"] = device.ID
	appServiceLabels["volume-id"] = volume.ID

	servicePorts := []apiv1.ServicePort{}
	for _, port := range app.OpenPorts {
		servicePorts = append(servicePorts, apiv1.ServicePort{
			Name:     fmt.Sprintf("svc-port-%d", port),
			Port:     int32(port),
			Protocol: apiv1.ProtocolTCP,
		})
	}

	// prepare selector based on stateful or stateless app
	selector := map[string]string{
		"app": adapter.GetAppDeploymentName(appRun.ID),
	}
	if app.Stateful {
		selector = map[string]string{
			"app": adapter.GetAppStatefulSetName(appRun.ID),
		}
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appServiceName,
			Labels:    appServiceLabels,
			Namespace: appServiceNamespace,
		},
		Spec: apiv1.ServiceSpec{
			Selector: selector,
			Type:     apiv1.ServiceTypeNodePort,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(volumeNamespace)
	service, err := serviceClient.Get(ctx, appServiceName, metav1.GetOptions{})
	// does not exist -> cannot update
	if err != nil {
		return err
	} else {
		service.ObjectMeta.Labels = appServiceLabels
		service.Spec.Ports = servicePorts
		service.Spec.Selector = selector

		// exist -> update
		_, updateErr := serviceClient.Update(ctx, service, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteAppService(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteAppService",
	})

	logger.Debug("received deleteAppService()")

	appServiceName := adapter.GetAppServiceName(appRunID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(appServiceNamespace)
	err := serviceClient.Delete(ctx, appServiceName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createAppIngress(app *types.App, appRun *types.AppRun) error {

	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "createAppIngress",
	})

	logger.Debug("received createAppIngress()")

	appIngressName := adapter.GetAppIngressName(appRun.ID)
	appIngressLabels := adapter.getAppIngressLabels(appRun)

	pathPrefix := networkingv1.PathTypePrefix

	serviceBackendPort := 0
	if len(app.OpenPorts) > 0 {
		serviceBackendPort = app.OpenPorts[0]
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appIngressName,
			Labels:    appIngressLabels,
			Namespace: appIngressNamespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                       "nginx",
				"nginx.ingress.kubernetes.io/proxy-connect-timeout": "150",
				"nginx.ingress.kubernetes.io/proxy-read-timeout":    "150",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     adapter.GetAppIngressPath(appRun.ID),
									PathType: &pathPrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: adapter.GetAppServiceName(appRun.ID),
											Port: networkingv1.ServiceBackendPort{
												Number: int32(serviceBackendPort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	ingressClient := adapter.clientSet.NetworkingV1().Ingresses(volumeNamespace)
	_, err := ingressClient.Get(ctx, ingress.GetName(), metav1.GetOptions{})

	if err != nil {
		// does not exist -> create
		_, createErr := ingressClient.Create(ctx, ingress, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := ingressClient.Update(ctx, ingress, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) updateAppIngress(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {

	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "updateAppIngress",
	})

	logger.Debug("received updateAppIngress()")

	appIngressName := adapter.GetAppIngressName(appRun.ID)
	appIngressLabels := adapter.getAppIngressLabels(appRun)
	// update labels
	appIngressLabels["app-id"] = app.ID
	appIngressLabels["device-id"] = device.ID
	appIngressLabels["volume-id"] = volume.ID

	pathPrefix := networkingv1.PathTypePrefix

	serviceBackendPort := 0
	if len(app.OpenPorts) > 0 {
		serviceBackendPort = app.OpenPorts[0]
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appIngressName,
			Labels:    appIngressLabels,
			Namespace: appIngressNamespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                       "nginx",
				"nginx.ingress.kubernetes.io/proxy-connect-timeout": "150",
				"nginx.ingress.kubernetes.io/proxy-read-timeout":    "150",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     adapter.GetAppIngressPath(appRun.ID),
									PathType: &pathPrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: adapter.GetAppServiceName(appRun.ID),
											Port: networkingv1.ServiceBackendPort{
												Number: int32(serviceBackendPort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	ingressClient := adapter.clientSet.NetworkingV1().Ingresses(volumeNamespace)
	_, err := ingressClient.Get(ctx, ingress.GetName(), metav1.GetOptions{})
	// does not exist -> cannot update
	if err != nil {
		return err
	} else {
		// exist -> update
		_, updateErr := ingressClient.Update(ctx, ingress, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteAppIngress(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteAppIngress",
	})

	logger.Debug("received deleteAppIngress()")

	appIngressName := adapter.GetAppIngressName(appRunID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	ingressClient := adapter.clientSet.NetworkingV1().Ingresses(appIngressNamespace)
	err := ingressClient.Delete(ctx, appIngressName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) CreateApp(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateApp",
	})

	logger.Debug("received CreateApp()")

	var err error

	if app.Stateful {
		err = adapter.createAppStatefulSet(device, volume, app, appRun)

	} else {
		err = adapter.createAppDeployment(device, volume, app, appRun)
	}
	if err != nil {
		return err
	}

	err = adapter.createAppService(app, appRun)
	if err != nil {
		return err
	}

	err = adapter.createAppIngress(app, appRun)
	if err != nil {
		panic(err)
	}

	return nil
}

func (adapter *K8SAdapter) UpdateAppRun(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "UpdateAppRun",
	})

	logger.Debug("received UpdateAppRun()")

	var err error

	if app.Stateful {
		err = adapter.updateAppStatefulSet(device, volume, app, appRun)
	} else {
		err = adapter.updateAppDeployment(device, volume, app, appRun)
	}

	if err != nil {
		return err
	}

	err = adapter.updateAppService(device, volume, app, appRun)
	if err != nil {
		return err
	}

	// err = adapter.updateAppIngress(device, volume, app, appRun)
	// if err != nil {
	// 	panic(err)
	// }

	return nil
}

func (adapter *K8SAdapter) DeleteApp(appRunID string, stateful bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "DeleteApp",
	})

	logger.Debug("received DeleteApp()")

	err := adapter.deleteAppIngress(appRunID)
	if err != nil {
		return err
	}

	err = adapter.deleteAppService(appRunID)
	if err != nil {
		return err
	}

	if stateful {
		err = adapter.deleteAppStatefulSet(appRunID)
	} else {
		err = adapter.deleteAppDeployment(appRunID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) EnsureDeleteApp(appRunID string) {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "EnsureDeleteApp",
	})

	logger.Debug("received EnsureDeleteApp()")

	adapter.deleteAppIngress(appRunID)
	adapter.deleteAppService(appRunID)
	adapter.deleteAppDeployment(appRunID)
}

