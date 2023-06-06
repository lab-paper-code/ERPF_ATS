package k8s

import (
	"fmt"
	"context"
	"k8s_old_ref"
		"context"
	"k8s_old_ref"
	
	"github.com/lab-paper-code/ksv/volume-service/types"
	"k8s.io/client-go/kubernetes"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/client-go/kubernetes"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"

)

const (
	appDeploymentNamePrefix string = "app"
	appDeploymentNamespace  string = objectNamespace
	appServiceNamePrefix    string = "app"
	appServiceNamespace     string = objectNamespace
	appIngressNamePrefix    string = "app"
	appIngressNamespace     string = objectNamespace

	appContainerVolumeName  string = "app-storage"
	appContainerPVMountPath string = "/uploads"
)

func (adapter *K8SAdapter) GetAppDeploymentName(appRunID string) string {
	return fmt.Sprintf("%s_%s", appDeploymentNamePrefix, appRunID)
}

func (adapter *K8SAdapter) GetAppServiceName(appRunID string) string {
	return fmt.Sprintf("%s_%s", appServiceNamePrefix, appRunID)
}

func (adapter *K8SAdapter) GetAppIngressName(appRunID string) string {
	return fmt.Sprintf("%s_%s", appIngressNamePrefix, appRunID)
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

func (adapter *K8SAdapter) CreateApp(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateApp",
	})

	logger.Debug("received CreateApp()")

	// TODO: Implement
	/*
		err := adapter.createAppDeployment(device, volume, app, appRun)
		if err != nil {
			return err
		}

		err = adapter.createAppService(appRun)
		if err != nil {
			return err
		}

		err = adapter.createAppIngress(appRun)
		if err != nil {
			panic(err)
		}
	*/

	return nil
}

func (adapter *K8SAdapter) DeleteApp(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "DeleteApp",
	})

	logger.Debug("received DeleteApp()")

	/*
		err := adapter.deleteAppIngress(appRunID)
		if err != nil {
			return err
		}

		err = adapter.deleteAppService(appRunID)
		if err != nil {
			return err
		}

		err = adapter.deleteAppDeployment(appRunID)
		if err != nil {
			return err
		}
	*/

	return nil
}

// App deployment example
/*
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pod1-app #변경
  namespace: ksv

spec:
  replicas: 1
  selector:
    matchLabels:
      app: pod1-app #변경
  revisionHistoryLimit: 2
  template:
    metadata:
      labels:
        app: pod1-app #변경
    spec:
      containers:
        - name: app-image
          #image: yechae/ksv-app:v3
          image: yechae/kube-flask:v4
          imagePullPolicy: IfNotPresent
          ports:
          - containerPort: 5000
          volumeMounts:
          - mountPath: "/mnt"
            name: volumes
          resources:
            requests:
              cpu: "250m"
            limits:
              cpu: "500m"

      volumes:
      - name: volumes
        persistentVolumeClaim:
          claimName: pod1-pvc #변경
      restartPolicy: Always
*/


// CreateAppDeploy creates a App deploy for the given volumeID
func (client *K8sClient) CreateAppDeploy(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateAppDeploy",
	})

	logger.Debugf("Creating a App Deploy for user %s, volume id %s", username, volumeID)

	deployAppName := client.getDeployAppName(volumeID)
	deployReplicas := int32(1)

	claim := &appsv1.Deployment{ // Deployment enables declarative updates for Pods and ReplicaSets.
		ObjectMeta: metav1.ObjectMeta{
			Name:	deployAppName,
			Labels:	client.getDeployLabels(username, volumeID),
			Namespace:	client.getDeployNamespace(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deployAppName,
				},
			},
			Template: corev1.PodTemplateSpec{ // PodTemplateSpec describes the data a pod should have when created from a template
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deployAppName,
					},
				},
			Spec: corev1.PodSpec{ // PodSpec is a description of a pod.
				Containers: []corev1.Container{ // A single application container that you want to run within a pod.
					{
						Name: "app-image",
						Image: "yechae/ksv-app:v4",
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{ // ContainerPort represents a network port in a single container.
							{
								ContainerPort: 5000,
							},
						},
						// Resources: corev1.ResourceRequirements{ // ResourceRequirements describes the compute resource requirements.
						// 	Requests: map[string]string{
						// 		cpu: "250m",

						// 	},
						// 	Limits: map[string]string{
						// 		cpu: "500m",
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{ // VolumeMount describes a mounting of a Volume within a container.
							{
								MountPath: "/mnt",
								Name: "volumes",
							},
						},
					},//Containers
					},//Containers
				Volumes: []corev1.Volume{ // Volume represents a named volume in a pod that may be accessed by any container in the pod.
					{
						Name: "volumes",
						VolumeSource: corev1.VolumeSource{ // Represents the source of a volume to mount.
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: client.getPVCName(volumeID),
							},

						},
					},
				},
				RestartPolicy: "Always",
				},//spec
			},
			},
		}


	deployclient := client.clientSet.AppsV1().Deployments(client.getDeployNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := deployclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = deployclient.Create(ctx, claim, metav1.CreateOptions{}) // CreateOptions may be provided when creating an API object.
		if err != nil {
			print(err,"\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a App Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a App Deploy for user %s, volume id %s", username, volumeID)
	} else {
		_, err = deployclient.Update(ctx, claim, metav1.UpdateOptions{}) // UpdateOptions may be provided when updating an API object.
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a App Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a App Deploy for user %s, volume id %s", username, volumeID)
	}

	return nil
}

	//make App ingress
	err = k8sClient.CreateAppIngress(input.Username, volumeID)
	if err != nil {
		panic(err)
	}

	err = k8sClient.WaitPodRun3(input.Username, volumeID)
	if err != nil {
		panic(err)
	}

	logger.Infof("All pods in podname=\"%s\" are running!", volumeID)

	execCommand := "sed -i -e 's#Alias /uploads \"/uploads\"#Alias /" + volumeID + "/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf"
	//change webdav path using volumeID
	err = k8sClient.ExecInPod("vd", volumeID, execCommand)
	if err != nil {
		panic(err)
	}

	execCommand = "/usr/sbin/httpd -k restart"
	err = k8sClient.ExecInPod("vd", volumeID, execCommand)
	if err != nil {
		panic(err)
	}

	//TODO:  k8s resource들 생성한 후
	//1. webdav pod으로 exec 명령어로 sed -i -e 's#Alias /uploads \"/uploads\"#Alias /<volumeID>/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf 명령어 실행
	//2. app pod으로 http://ip:60000/hello_flask?ip=<ip> 해서 dom ip 알려주기

	type Output struct {
		Mount  string       `json:mountPath`
		Device types.Device `json: device`
	}

	Mount := "http://155.230.36.27/" + volumeID + "/uploads"
	// 지금과 다름
	device = types.Device{
		IP:       input.IP,
		ID:       volumeID,
		Username: input.Username,
		Password: input.Password,
		Storage:  input.Storage,
	}

	output := Output{
		Mount:  Mount,
		Device: device,
	}
