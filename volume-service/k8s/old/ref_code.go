package k8s

import (
	"context"
	"fmt"
	"os"

	"github.com/rook/rook/pkg/client/clientset/versioned/scheme"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// WaitPodRun3
func (client *K8sClient) WaitPodRun3(username string, volumeID string) error {
	// pod, err := client.PodClient().Get(podName, metav1.GetOptions{})
	// if err!= nil {
	// 	panic(err)
	// }

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	namespace := client.getDeployNamespace()
	podName := client.getPodName(volumeID)

	pod, err := client.clientSet.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	ctx := context.Background()
	watcher, err := coreclient.Pods(namespace).Watch(
		ctx,
		metav1.SingleObject(pod.ObjectMeta),
	)
	if err != nil {
		return err
	}

	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			pod := event.Object.(*corev1.Pod)

			if pod.Status.Phase == corev1.PodRunning {
				fmt.Printf("The POD is runnging")
				return nil
			}

		case <-ctx.Done():
			fmt.Printf("Exit from waitPodRunning for POD \"%s\" because the context is done", volumeID)
			return nil
		}
	}

	return nil
}

// getPodName
func (client *K8sClient) getPodName(volumeID string) string {

	podLabel := client.getDeployWebdavName(volumeID)
	pods, err := client.clientSet.CoreV1().Pods("vd").List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=" + podLabel,
	})

	if err != nil {
		panic(err)
	}

	var podName string
	for _, pod := range pods.Items {
		podName = pod.Name
	}
	fmt.Println(podName)

	return podName
}

// K8sClient struct
type K8sClient struct {
	kubeConfigPath string
	clientSet      *kubernetes.Clientset // Clientset contains the clients for groups
}

// getDeployAppName
func (client *K8sClient) getDeployAppName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, deployAppSuffix)
}

// getDeployWebdavName
func (client *K8sClient) getDeployWebdavName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, deployWebdavSuffix)
}

// getDeployLabels
func (client *K8sClient) getDeployLabels(username string, volumeID string) map[string]string {
	return map[string]string{
		"username":  username,
		"volume-id": volumeID,
	}
}

// getDeployNamespace
func (client *K8sClient) getDeployNamespace() string {
	return volumeNamespace
}

// getPVCName
func (client *K8sClient) getPVCName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, pvcSuffix)
}

// ExecInPod
func (client *K8sClient) ExecInPod(namespace string, volumeID string, command string) error {

	// pod, err := client.PodClient().Get(podName, metav1.GetOptions{})
	// if err!= nil {
	// 	panic(err)
	// }

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	//coerv1client change to client
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	// Prepare the API URL used to execute another process within the Pod.  In
	// this case, we'll run a remote shell.

	podLabel := client.getDeployWebdavName(volumeID)
	pods, err := client.clientSet.CoreV1().Pods("vd").List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=" + podLabel,
	})

	if err != nil {
		panic(err)
	}

	var podName string
	//var podObj corev1.Pod
	for _, pod := range pods.Items {
		podName = pod.Name
		//podObj = pod
	}

	execCommand := []string{
		"sh",
		"-c",
		command,
	}

	fmt.Println(execCommand)
	req := coreclient.RESTClient().
		Post().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			// Container: podObj.Spec.Containers[0].Name,
			Command: execCommand,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		print(err)
	}

	return nil
}

func (client *K8sClient) getVolumeNamespace() string {
	return volumeNamespace
}

// Service
func (client *K8sClient) getWebdavSVCName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, svcWebdavSuffix)
}

// Ingress
func (client *K8sClient) CreateAppIngress(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateAppIngress",
	})

	logger.Debugf("Creating a App Ingress for user %s, volume id %s", username, volumeID)
	pathPrefix := networkingv1.PathTypePrefix

	claim := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      client.getAppIngressName(volumeID),
			Namespace: client.getIngressNamespace(),
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
									Path:     client.getAppIngressPath(volumeID),
									PathType: &pathPrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: client.getAppSVCName(volumeID),
											Port: networkingv1.ServiceBackendPort{
												Number: 60000,
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

	// claim := &extensionsv1beta1.Ingress{
	// 	ApiVersion: "networking.k8s.io/v1",
	// 	Kind: "Ingress",
	// 	ObjectMeta: metav1.ObjectMetat{
	// 		Name: client.getAppIngressName(volumeID),
	// 		Namespace: client.getIngressNamespace(),
	// 		Annotations: {
	// 			kubernetes.io/ingress.class: "nginx",
	// 			nginx.ingress.kubernetes.io/proxy-connect-timeout: "150",
	// 			nginx.ingress.kubernetes.io/proxy-read-timeout: "150",

	// 		},
	// 	},
	// 	Spec: []extensionsv1beta1.IngressRule{
	// 		Http: []extensionsv1beta1.HttpIngressPath{
	// 			Path: client.getAppIngressPath(volumeID),
	// 			Backend: extensionsv1beta1.IngressBackend{
	// 				ServiceName: client.getAppSVCName(volumeID),
	// 				ServicePort: map[string]String{
	// 					60000
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	appIngclient := client.clientSet.NetworkingV1().Ingresses(client.getVolumeNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := appIngclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = appIngclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err, "\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a appSVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a appSVC for user %s, volume id %s", username, volumeID)
	} else {
		_, err = appIngclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a appSVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a appSVC for user %s, volume id %s", username, volumeID)
	}

	return nil
}
