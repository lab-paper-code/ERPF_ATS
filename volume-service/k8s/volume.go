package k8s

import (
	"context"
	"fmt"

	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	resourcev1 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	volumeNamePrefix      string = "pv"
	volumeClaimNamePrefix string = "pvc"
	volumeNamespace       string = objectNamespace
	storageClassName      string = "rook-cephfs"
)

func (adapter *K8SAdapter) GetStorageClassName() string {
	return storageClassName
}

func (adapter *K8SAdapter) GetVolumeName(device *types.Device) string {
	return fmt.Sprintf("%s_%s", volumeNamePrefix, device.ID)
}

func (adapter *K8SAdapter) GetVolumeClaimName(device *types.Device) string {
	return fmt.Sprintf("%s_%s", volumeClaimNamePrefix, device.ID)
}

func (adapter *K8SAdapter) GetVolumeLabels(device *types.Device) map[string]string {
	labels := map[string]string{}
	labels["volume-name"] = adapter.GetVolumeName(device)
	labels["device-id"] = device.ID
	return labels
}

func (adapter *K8SAdapter) getCephCSIPersistentVolumeSource(device *types.Device) *apiv1.CSIPersistentVolumeSource {
	// TODO: Implement this
	return nil
}

func (adapter *K8SAdapter) CreatePV(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreatePV",
	})

	logger.Debug("received CreatePV()")

	volumeName := adapter.GetVolumeName(device)
	volumeLabels := adapter.GetVolumeLabels(device)

	volumeSize := resourcev1.Quantity{
		Format: resourcev1.BinarySI,
	}
	volumeSize.Set(device.VolumeSize)

	volmode := apiv1.PersistentVolumeFilesystem

	cephCSIPersistentVolumeSource := adapter.getCephCSIPersistentVolumeSource(device)

	pv := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      volumeName,
			Labels:    volumeLabels,
			Namespace: volumeNamespace,
		},
		Spec: apiv1.PersistentVolumeSpec{
			Capacity: apiv1.ResourceList{
				apiv1.ResourceStorage: volumeSize,
			},
			VolumeMode: &volmode,
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteMany,
			},
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimRetain,
			StorageClassName:              adapter.GetStorageClassName(),
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				CSI: cephCSIPersistentVolumeSource,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	pvclient := adapter.clientSet.CoreV1().PersistentVolumes()
	_, err := pvclient.Get(ctx, pv.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := pvclient.Create(ctx, pv, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := pvclient.Update(ctx, pv, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) CreatePVC(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreatePVC",
	})

	logger.Debug("received CreatePVC()")

	volumeName := adapter.GetVolumeName(device)
	volumeClaimName := adapter.GetVolumeClaimName(device)
	volumeLabels := adapter.GetVolumeLabels(device)

	// we request very small size volume to match any pv available
	// TODO: double check if this pvc can bind to pv
	volumeSize := resourcev1.Quantity{
		Format: resourcev1.BinarySI,
	}
	volumeSize.Set(1024)

	storageClassName := adapter.GetStorageClassName()

	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      volumeClaimName,
			Labels:    volumeLabels,
			Namespace: volumeNamespace,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			VolumeName: volumeName,
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteMany,
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: volumeSize,
				},
			},
			StorageClassName: &storageClassName,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"volume-name": volumeName,
				},
			},
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
