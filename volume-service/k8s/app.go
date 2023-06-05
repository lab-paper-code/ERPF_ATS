package k8s

import (
	"fmt"

	"github.com/lab-paper-code/ksv/volume-service/types"
)

const (
	appDeploymentNamePrefix string = objectNamespace
)

func (adapter *K8SAdapter) GetAppDeploymentName(device *types.Device) string {
	return fmt.Sprintf("%s_%s", appDeploymentNamePrefix, device.ID)
}

//APP
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

/*
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

	claim := &appsv1.Deployment{
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
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deployAppName,
					},
				},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "app-image",
						Image: "yechae/ksv-app:v4",
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 5000,
							},
						},
						// Resources: corev1.ResourceRequirements{
						// 	Requests: map[string]string{
						// 		cpu: "250m",

						// 	},
						// 	Limits: map[string]string{
						// 		cpu: "500m",
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/mnt",
								Name: "volumes",
							},
						},
					},//Continers
					},//Continers
				Volumes: []corev1.Volume{
					{
						Name: "volumes",
						VolumeSource: corev1.VolumeSource{
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
		_, err = deployclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err,"\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a App Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a App Deploy for user %s, volume id %s", username, volumeID)
	} else {
		_, err = deployclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a App Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a App Deploy for user %s, volume id %s", username, volumeID)
	}

	return nil
}
*/
/*
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
	//1. webdav pod으로 exec 명령어로 sed -i -e 's#Alias /uploads \"/uploads\"#Alias /<volumdID>/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf 명령어 실행
	//2. app pod으로 http://ip:60000/hello_flask?ip=<ip> 해서 dom ip 알려주기

	type Output struct {
		Mount  string       `json:mountPath`
		Device types.Device `json: device`
	}

	Mount := "http://155.230.36.27/" + volumeID + "/uploads"

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
*/
