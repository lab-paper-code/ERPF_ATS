package logic

import (
	"awds/types"
	"fmt"
	"math"
	"strconv"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var batchSize int = 10 // need to change later

func (logic *Logic) ListJobs() ([]types.Job, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListJobs",
	})

	logger.Debug("received ListJobs()")

	return logic.dbAdapter.ListJobs()
}

func (logic *Logic) GetJob(jobID string) (types.Job, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetJob",
	})

	logger.Debug("received GetJob()")

	return logic.dbAdapter.GetJob(jobID)
}

func (logic *Logic) CreateJob(job *types.Job) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateJob",
	})

	logger.Debug("received CreateJob()")

	return logic.dbAdapter.InsertJob(job)
}

func (logic *Logic) UpdateJobDevice (jobID string, deviceID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateJobDevice",
	})

	logger.Debug("received UpdateJobDevice()")

	return logic.dbAdapter.UpdateJobDevice(jobID, deviceID)
}

func (logic *Logic) UpdateJobPod(jobID string, podID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateJobPod",
	})

	logger.Debug("received UpdateJobPod()")

	return logic.dbAdapter.UpdateJobPod(jobID, podID)
}

func (logic *Logic) UpdateInputSize(jobID string, inputSize int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateInputSize",
	})

	logger.Debug("received UpdateInputSize()")

	return logic.dbAdapter.UpdateInputSize(jobID, inputSize)
}

func (logic *Logic) UpdatePartitionRate(jobID string, partitionRate float64) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePartitionRate",
	})

	logger.Debug("received UpdatePartitionRate()")

	return logic.dbAdapter.UpdatePartitionRate(jobID, partitionRate)
}

func (logic *Logic) UpdateDeviceStartIndex(jobID string, deviceStartIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceStartIndex",
	})

	logger.Debug("received UpdateDeviceStartIndex()")

	return logic.dbAdapter.UpdateDeviceStartIndex(jobID, deviceStartIndex)
}


func (logic *Logic) UpdateDeviceEndIndex(jobID string, deviceEndIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceEndIndex",
	})

	logger.Debug("received UpdateDeviceEndIndex()")

	return logic.dbAdapter.UpdateDeviceEndIndex(jobID, deviceEndIndex)
}


func (logic *Logic) UpdatePodStartIndex(jobID string, podStartIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePodStartIndex",
	})

	logger.Debug("received UpdatePodStartIndex()")

	return logic.dbAdapter.UpdatePodStartIndex(jobID, podStartIndex)
}



func (logic *Logic) UpdatePodEndIndex(jobID string, podEndIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePodEndIndex",
	})

	logger.Debug("received UpdatePodEndIndex()")

	return logic.dbAdapter.UpdatePodEndIndex(jobID, podEndIndex)
}

func (logic *Logic) DeleteJob(jobID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeleteJob",
	})

	logger.Debug("received DeleteJob()")

	return logic.dbAdapter.DeleteJob(jobID)
}

func (logic *Logic) GetFullEndpoint(endpoint string, startIdx int, endIdx int) string {
	return fmt.Sprintf("%s/%d-%d", endpoint, startIdx, endIdx)
}

func (logic *Logic) Precompute(jobID string, deviceEndpoint string, podEndpoint string, inputSize int) error {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Precompute",
	})

	logger.Debug("received Precompute()")

	var deviceResponse map[string]interface{}
	var podResponse map[string]interface{}
	
	startIdx := 0
	endIdx := 1	

	precomputeSize := int( 0.001 * float64(inputSize) )
	// set endIdx if precomputeSize is bigger than 1
	if precomputeSize >= 1 {
		endIdx = precomputeSize
	}
	
	deviceFullEndpoint := logic.GetFullEndpoint(deviceEndpoint, startIdx, endIdx)

	podFullEndpoint := logic.GetFullEndpoint(podEndpoint, startIdx, endIdx)

	fmt.Println(deviceFullEndpoint)
	fmt.Println(podFullEndpoint)

	client := resty.New()
	deviceResultChan := make(chan float64, 1)
	podResultChan := make(chan float64, 1)
	errChan := make(chan error, 2)

	go func() {
		_, deviceRequestErr := client.R().
			SetResult(&deviceResponse).
			Get(deviceFullEndpoint)
		if deviceRequestErr != nil {
			errChan <- deviceRequestErr
			return
		}

		deviceResultStr, ok := deviceResponse["result"].(string)
		if !ok {
			errChan <- fmt.Errorf("device response 'result' is not a string")
			return
		}

		deviceResult, err := strconv.ParseFloat(deviceResultStr, 64)
		if err != nil {
			errChan <- fmt.Errorf("failed to parse device result to float64: %v", err)
			return
		}
		deviceResultChan <- deviceResult
	}()

	go func() {
		_, podRequestErr := client.R().
		SetResult(&podResponse).
		Get(podFullEndpoint)	  
	if podRequestErr != nil {
		errChan <- podRequestErr
		return
	}

	podResultStr, ok := podResponse["result"].(string)
	if !ok {
		errChan <- fmt.Errorf("pod response 'result' is not a string")
		return
	}

	podResult, err := strconv.ParseFloat(podResultStr, 64)
	if err != nil {
		errChan <- fmt.Errorf("failed to parse pod result to float64: %v", err)
		return
	}
	
	podResultChan <- podResult
	
	}()
	
	var deviceResult, podResult float64
	for i := 0; i < 2; i++ {
		select {
		case deviceResult = <- deviceResultChan:
		case podResult = <- podResultChan:
		case err := <- errChan:
			return err
		}
	}

	// set StartIdx, EndIdx based on precomputation result
	partitionRate := math.Round(podResult / (deviceResult + podResult) * 100) / 100
	fmt.Println("deviceResult", deviceResult)
	fmt.Println("podResult", podResult)
	fmt.Println("partitionRate", partitionRate)

	deviceEndIdx := int( partitionRate * float64(batchSize) )
	fmt.Println("Before save, deviceEndIdx", deviceEndIdx)

	// if job finishes in single distribution, set batchSize to inputSize
	if batchSize >= inputSize {
		batchSize = inputSize		
	}

	// update device start, end index
	err := logic.dbAdapter.UpdateDeviceStartIndex(jobID, 0)
	if err != nil {
		return err
	}

	err = logic.dbAdapter.UpdateDeviceEndIndex(jobID, deviceEndIdx)
	if err != nil {
		return err
	}
	
	// update pod start, end index
	err = logic.dbAdapter.UpdatePodStartIndex(jobID, deviceEndIdx)
	if err != nil {
		return err
	}
	
	err = logic.dbAdapter.UpdatePodEndIndex(jobID, batchSize)
	if err != nil {
		return err
	}
	
	err = logic.dbAdapter.UpdatePartitionRate(jobID, partitionRate)
	return nil
}

func (logic *Logic) ComputeDevice(jobID string, deviceEndpoint string, startIdx int, endIdx int) (float64, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ComputeDevice",
	})

	logger.Debugf("received ComputeDevice()")

	
	var deviceResponse map[string]interface{}

	client := resty.New()
	// request to device	
	deviceFullEndpoint := logic.GetFullEndpoint(deviceEndpoint, startIdx, endIdx)

	fmt.Println("compute device full endpoint", deviceFullEndpoint)

	_, err := client.R().
	SetResult(&deviceResponse).
	Get(deviceFullEndpoint)
	if err != nil {
		return -1, err
	}

	deviceResultStr, ok := deviceResponse["result"].(string)
	if !ok {
		return -1, fmt.Errorf("device response 'result' is not a string")
	}

	deviceResult, err := strconv.ParseFloat(deviceResultStr, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse result to float64: %v", err)
	}

	fmt.Println("device", deviceResult)
	
	return deviceResult, nil
}

func (logic *Logic) ComputePod(jobID string, podEndpoint string, startIdx int, endIdx int) (float64, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ComputePod",
	})

	logger.Debugf("received ComputePod()")

	
	var podResponse map[string]interface{}

	client := resty.New()
		// request to device	
	podFullEndpoint := logic.GetFullEndpoint(podEndpoint, startIdx, endIdx)

	fmt.Println("compute pod full endpoint", podFullEndpoint)

	_, err := client.R().
	SetResult(&podResponse).
	Get(podFullEndpoint)
	if err != nil {
		return -1, err
	}

	podResultStr, ok := podResponse["result"].(string)
	if !ok {
		return -1, fmt.Errorf("pod response 'result' is not a string")
	}

	podResult, err := strconv.ParseFloat(podResultStr, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse result to float64: %v", err)
	}
	// request to pod
	
	fmt.Println("pod", podResult)

	return podResult, nil

	}

func (logic *Logic) ScheduleJob(jobID string) error {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ScheduleJob",
	})

	logger.Debug("received ScheduleJob()")

	job, err := logic.dbAdapter.GetJob(jobID)
	if err != nil {
		return err
	}

	device, err := logic.dbAdapter.GetDevice(job.DeviceID)
	if err != nil {
		return err
	}

	pod, err := logic.dbAdapter.GetPod(job.PodID)
	if err != nil {
		return err
	}

	// precompute
	err = logic.Precompute(jobID, device.Endpoint, pod.Endpoint, job.InputSize)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("precomputation ended...")
	// compute til end
	for {
		job, err = logic.dbAdapter.GetJob(jobID)
		if err != nil {
			fmt.Errorf("%v", err)
			return err
		}
		
		// break when job completes
		if job.Completed {
			break
		}
		
		fmt.Println("job.DeviceStartIndex in compute", job.DeviceStartIndex)
		fmt.Println("job.DeviceEndIndex in compute", job.DeviceEndIndex)
		
		// use channels to handle concurrent Compute tasks
		deviceResultsChan := make(chan float64, 1)
		podResultsChan := make(chan float64, 1)
		errChan := make(chan error, 2)

		// concurrently run Compute for device
		go func() {
			result, err := logic.ComputeDevice(jobID, device.Endpoint, job.DeviceStartIndex, job.DeviceEndIndex)
			if err != nil {
				errChan <- err
				return
			}
			deviceResultsChan <- result
		}()

		// concurrently run Compute for pod
		go func() {
			result, err := logic.ComputePod(jobID, pod.Endpoint, job.PodStartIndex, job.PodEndIndex)
			if err != nil {
				errChan <- err
				return
			}
			podResultsChan <- result
		}()

		var devResult, podResult float64
		for i :=0; i < 2; i++{
			select {
			case result := <-deviceResultsChan:
				devResult = result
			case result := <-podResultsChan:
				podResult = result
			case err := <-errChan:
				return err
			}
		}

		partitionRate := math.Round(podResult / (devResult + podResult) * 100) / 100

		if podResult < float64(1){
			// pod didn't do work
			// get previousPartitionRate
			partitionRate = job.PartitionRate	
			
		}

		fmt.Println("partitionRate", partitionRate)

		// update device start, end index
		// Assuming batchSize is 100 and partitionRate has been calculated
		deviceWork := int(float64(batchSize) * partitionRate)
		podWork := batchSize - deviceWork // Ensure total work does not exceed batchSize
		fmt.Println("deviceWork", deviceWork)
		fmt.Println("podWork", podWork)

		newDeviceStartIndex := job.PodEndIndex
		newDeviceEndIndex := job.PodEndIndex + deviceWork

		// Update device start and end index
		if job.PodEndIndex >= job.InputSize {
			newDeviceStartIndex = job.InputSize
		}

		if newDeviceEndIndex >= job.InputSize {
			newDeviceEndIndex = job.InputSize // Adjust to not exceed InputSize
		}

		err = logic.dbAdapter.UpdateDeviceStartIndex(jobID, newDeviceStartIndex)
		if err != nil {
			return err
		}
		err = logic.dbAdapter.UpdateDeviceEndIndex(jobID, newDeviceEndIndex)
		if err != nil {
			return err
		}

		newPodStartIndex := job.PodEndIndex + deviceWork
		newPodEndIndex := job.PodEndIndex + deviceWork + podWork

		// Update pod start and end index
		if newPodStartIndex > job.InputSize {
			newPodStartIndex = job.InputSize // Adjust to not exceed InputSize
		}
		if newPodEndIndex > job.InputSize {
			newPodEndIndex = job.InputSize // Adjust to not exceed InputSize
		}
		err = logic.dbAdapter.UpdatePodStartIndex(jobID, newPodStartIndex)
		if err != nil {
			return err
		}
		err = logic.dbAdapter.UpdatePodEndIndex(jobID, newPodEndIndex)
		if err != nil {
			return err
		}
		
		// job completed
		// TODO2: need to test this code
		if ( job.DeviceStartIndex == job.InputSize ) || ( job.PodStartIndex == job.InputSize ){
			err = logic.dbAdapter.UpdateJobCompleted(jobID, true)
			if err != nil {
				return err
			}
		}		
	}

	return nil
}

	// 먼저 끝나면, 끝난 장치에 요청 다시 보내야
	// 다음 배치 가져와서 실행
	// 파드 먼저 끝남, 디바이스 예측 시간이랑 비슷하게 끝나게 설정
	// 파드에서 끝난 시간 알 수 있음
	// 파드에서 * 2로 처리 -> 나중에 변경
	// 디바이스에서는 한 파티션 끝날 때 파드가 두 번째 게 돌고 있을 것 같으므로
	// 파드에서 2번째 것 끝날 시간이랑 거의 겹치게 다음 배치 데이터를 넣어줌
	// 마지막으로 끝난 애 기준으로, 늦은 애가 끝났을 때 빠른 애가 계속 돌고 있으므로
	// 거기에 맞춰서 빠른 애가 몇 번 맞췄는지, 
	// 마지막 세트: 20 = 8 + 8 + 4, 마지막 배치 비율로


	// type Scheduler struct {
		// 	config *commons.Config
		// 	Job_list []*types.Job
		// }
		
		// type ScheduleJob interface {
		// 	getJob(string) (*types.Job)
		// 	Precompute(*types.Job, string, string) (float64, error)
		// 	Compute(*types.Job, string, string) (error)
		// 	// add method if necessary
		// }
		
		// // getJob returns Job from jobID
		// func (scheduler *Scheduler) GetJob (jobID string) (types.Job) {
		// 	var idx int
		// 	jobList := scheduler.Job_list
		
		// 	// using generic, supported after Go 1.21
		// 	// idx = slices.IndexFunc(scheduler.Job_list, func(j_ptr *types.Job) bool { return (*j_ptr).ID == jobID })
		
		// 	// using for loop, slower than Generic, use this for Go version lower than 1.21
		// 	for i := range jobList {
		// 		if jobList[i].ID == jobID {
		// 			idx = i
		// 			break
		// 		}
		// 	}
		
		// 	return *jobList[idx]
		// }

		
// func (logic *Logic) UnscheduleJob(jobID string) error {
// 	logger := log.WithFields(log.Fields{
// 		"package": "logic",
// 		"struct" : "Logic",
// 		"function" : "UnscheduleJob",
// 	})

// 	logger.Debug("received UnscheduleJob()")

// 	job, err := logic.dbAdapter.GetJob(jobID)
// 	if err != nil {
// 		return err
// 	}

// 	pod, err := logic.dbAdapter.GetPod(job.PodID)
// 	if err != nil {
// 		return err
// 	}

// 	for true {
// 		job, _ = logic.dbAdapter.GetJob(jobID)
// 		Compute()
		
// 		// break when job completes
// 		if job.Completed /* &&  device Job and pod Job completes  */ {
// 			break
// 		}
// 	}
// }