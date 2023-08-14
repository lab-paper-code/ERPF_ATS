package k8s

import (
	"context"

	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	resourcev1 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	volumeClaimNamePrefix string = "pvc"
	volumeNamespace       string = objectNamespace
	storageClassName      string = "rook-cephfs"
)

func (adapter *K8SAdapter) GetStorageClassName() string {
	return storageClassName
}

func (adapter *K8SAdapter) GetVolumeClaimName(volumeID string) string { // changed to avoid kubernetes error
	return makeValidObjectName(volumeClaimNamePrefix, volumeID)
}

func (adapter *K8SAdapter) getVolumeLabels(volume *types.Volume) map[string]string {
	labels := map[string]string{}
	labels["volume-id"] = volume.ID
	labels["device-id"] = volume.DeviceID
	return labels
}

func (adapter *K8SAdapter) CreateVolume(volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateVolume",
	})

	logger.Debug("received CreateVolume()")

	volumeClaimName := adapter.GetVolumeClaimName(volume.ID)
	volumeLabels := adapter.getVolumeLabels(volume)

	volumeSize := resourcev1.Quantity{
		Format: resourcev1.BinarySI,
	}
	volumeSize.Set(volume.VolumeSize)

	storageClassName := adapter.GetStorageClassName()

	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      volumeClaimName,
			Labels:    volumeLabels,
			Namespace: volumeNamespace,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteMany,
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: volumeSize,
				},
			},
			StorageClassName: &storageClassName,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	pvcclient := adapter.clientSet.CoreV1().PersistentVolumeClaims(volumeNamespace)
	_, err := pvcclient.Get(ctx, pvc.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := pvcclient.Create(ctx, pvc, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := pvcclient.Update(ctx, pvc, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) ResizeVolume(volumeID string, size int64) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "ResizeVolume",
	})

	logger.Debug("received ResizeVolume()")

	volumeClaimName := adapter.GetVolumeClaimName(volumeID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	pvcclient := adapter.clientSet.CoreV1().PersistentVolumeClaims(volumeNamespace)
	pvc, err := pvcclient.Get(ctx, volumeClaimName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// TODO: double check if this works -> not working, figuring out how to fix this
	pvc.Spec.Resources.Requests.Storage().Set(size)

	_, updateErr := pvcclient.Update(ctx, pvc, metav1.UpdateOptions{})
	if updateErr != nil {
		return updateErr
	}

	return nil
}
