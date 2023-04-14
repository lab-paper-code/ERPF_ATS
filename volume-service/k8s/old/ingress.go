package k8s

// corev1 "k8s.io/api/core/v1"
// resourcev1 "k8s.io/apimachinery/pkg/api/resource"

//"k8s.io/api/networking/v1beta1"

/*
const (
	ingWebdavSuffix         string = "-webdav-ing"
	ingAppSuffix            string = "-app-ing"
	ingNamespace   			string = "vd"
	ingWebdavPathSuffix     string = "/"
	ingAppPathSuffix		string = "/app/"
)


apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pod1-ingress # 변경
  namespace: ksv
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "150"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "150"
spec:
  rules:
    - host:
      http:
        paths:
          - path: /pod1  # 변경 # volumeID 로
            pathType: Prefix
            backend:
              service:
                name: webdav-pod1-svc # 변경
                port:
                  number: 80


// getWebdavSvcName makes webdavIngress name
func (client *K8sClient) getWebdavIngressName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, ingWebdavSuffix)
}

// getAppSvcName makes appIngress name
func (client *K8sClient) getAppIngressName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, ingAppSuffix)
}

func (client *K8sClient) getIngressNamespace() string {
	return ingNamespace
}

func (client *K8sClient) getWebdavIngressPath(volumeID string) string {
	return fmt.Sprintf("%s%s", ingWebdavPathSuffix, volumeID)
}

func (client *K8sClient) getAppIngressPath(volumeID string) string {
	return fmt.Sprintf("%s%s", ingAppPathSuffix, volumeID)
}




// CreateIngress Webdav creates a pvc for the given volumeID
func (client *K8sClient) CreateWebdavIngress(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateWebdavIngress",
	})

	logger.Debugf("Creating a Webdav Ingress for user %s, volume id %s", username, volumeID)

	pathPrefix := networkingv1.PathTypePrefix

	claim := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: client.getWebdavIngressName(volumeID),
			Namespace: client.getIngressNamespace(),
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				"nginx.ingress.kubernetes.io/proxy-connect-timeout": "150",
    			"nginx.ingress.kubernetes.io/proxy-read-timeout": "150",

			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: client.getWebdavIngressPath(volumeID),
									PathType: &pathPrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: client.getWebdavSVCName(volumeID),
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
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
	//beta
	// claim := &v1beta1.Ingress{
	// 	ObjectMeta: metav1.ObjectMetat{
	// 		Name: client.getWebdavIngressName(volumeID),
	// 		Namespace: client.getIngressNamespace(),
	// 		Annotations: map[string]string{
	// 			"kubernetes.io/ingress.class": "nginx",
	// 			"nginx.ingress.kubernetes.io/proxy-connect-timeout": "150",
    // 			"nginx.ingress.kubernetes.io/proxy-read-timeout": "150",

	// 		},
	// 	},
	// 	Spec: v1beta1.InressSpec{
	// 		Rules: []v1beta1.IngressRule{
	// 			{
	// 				IgressRuleValue: v1beta1.IngressRuleValue{
	// 					HTTP: &v1beta1.HTTPIngressRuleValue{
	// 						Paths: []v1beta1.HTTPIngressPath{
	// 							{Path: client.getWebdavIngressPath(volumeID)},
	// 						},
	// 						Backend: v1beta1.IngressBackend{
	// 							ServiceName: client.getWebdavSVCName(volumeID),
	// 							ServicePort: map[string]interface{},
	// 						},
	// 					},
	// 				},
	// 			},

	// 		},
	// 	},
	// }

	webdavIngclient := client.clientSet.NetworkingV1().Ingresses(client.getVolumeNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := webdavIngclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = webdavIngclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err,"\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a Webdav ingress for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a Webdav ingress for user %s, volume id %s", username, volumeID)
	} else {
		_, err = webdavIngclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a Webdav ingress for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a Webdav ingress for user %s, volume id %s", username, volumeID)
	}

	return nil
}


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
			Name: client.getAppIngressName(volumeID),
			Namespace: client.getIngressNamespace(),
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				"nginx.ingress.kubernetes.io/proxy-connect-timeout": "150",
    			"nginx.ingress.kubernetes.io/proxy-read-timeout": "150",

			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: client.getAppIngressPath(volumeID),
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
			print(err,"\n")
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
*/
