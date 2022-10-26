package k8s

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	kubeConfigPath string
	clientSet      *kubernetes.Clientset
}

// NewK8sClient creates a new K8sClient
func NewK8sClient(kubeConfigPath string) (*K8sClient, error) {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"function": "NewK8sClient",
	})

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		kubeConfigPath: kubeConfigPath,
		clientSet:      clientset,
	}, nil
}
