package logic

import (
	"awds/types"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// map to record device info for calculating batch size
type deviceRecord map[string][]float64

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

func (logic *Logic) InsertJob(job *types.Job) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "InsertJob",
	})

	logger.Debug("received InsertJob()")

	return logic.dbAdapter.InsertJob(job)
}

func (logic *Logic) UpdateDeviceIDList(jobID string, deviceIDList []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceIDList",
	})

	deviceIDTemp := []string{}
	
	for _, deviceID := range deviceIDList{
		deviceIDTemp = append(deviceIDTemp, deviceID)
	}

	deviceIDListCSV := strings.Join(deviceIDTemp, ",") // make deviceID list into string

	logger.Debug("received UpdateJobDevice()")

	return logic.dbAdapter.UpdateDeviceIDList(jobID, deviceIDListCSV)
}

func (logic *Logic) UpdateEndIndex(jobID string, endIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateEndIndex",
	})

	logger.Debug("received UpdateEndIndex()")

	return logic.dbAdapter.UpdateEndIndex(jobID, endIndex)
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

// return first batchsize according to elapsedTime of adjustment stage
func (logic *Logic) determineDeviceType(elapsedTime float64) (float64, int) {
	diff := math.Abs(logic.config.PrecomputeReferenceLatencies[0] - elapsedTime)
	idx := 0
	var predictedTime float64

	for i, v := range(logic.config.PrecomputeReferenceLatencies) {
		if math.Abs(v - elapsedTime) < diff {
			idx = i
			diff = math.Abs(v - elapsedTime)
		}
	}
	averageTilt :=  (elapsedTime + logic.config.PrecomputeReferenceLatencies[idx]) / float64(logic.config.InitialBatchSize * 2)
	nextbatch := int(logic.config.PrecomputeReferenceLatencies[idx] / averageTilt)
	predictedTime = averageTilt * float64(nextbatch)


	if nextbatch > logic.config.MaxBatchSize{
		nextbatch = logic.config.MaxBatchSize
		predictedTime = elapsedTime
	}
	
	return predictedTime, nextbatch
}

// ywjang
func (logic *Logic) AdjustBatchSize(devIdQ *Queue, devRcdMap *deviceRecord, startIdx int, endIdx int) (int, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "AdjustBatchSize",
	})
	logger.Debug("received AdjustBatchSize()")

	deviceNum := len(*devIdQ)
	adjustBatchSize := logic.config.InitialBatchSize
	
	if (endIdx - startIdx) < adjustBatchSize{
		adjustBatchSize = (endIdx - startIdx)
	}

	var wg sync.WaitGroup
	wg.Add(deviceNum)
	errChan := make(chan error, 1)

	for idx, deviceID := range *devIdQ{
		device, err := logic.dbAdapter.GetDevice(deviceID)
		if err != nil {
			return 0, err
		}

		func(i int)() {
			defer wg.Done() 
			elapsedTime, _, err := logic.Compute(&device, startIdx, startIdx + adjustBatchSize)
			if err != nil {
				errChan <- err
				return 
			}
			nextPredictTime, nextBatchSize  := logic.determineDeviceType(elapsedTime)
			
			fmt.Println("AdjustBatchSize(deviceID, elapsedTime, adjustBatchSize, nextBatchSize, predictTime): ", device.ID, elapsedTime, adjustBatchSize, nextBatchSize, nextPredictTime, time.Now().Unix())
			(*devRcdMap)[deviceID][0] = nextPredictTime // current PredictTime
			(*devRcdMap)[deviceID][1] = elapsedTime // current elapsedTime
			(*devRcdMap)[deviceID][2] = float64(logic.config.InitialBatchSize) // current BatchSize
			(*devRcdMap)[deviceID][3] = float64(nextBatchSize)// nextBatchSize
			(*devRcdMap)[deviceID][4] += 1 // used to count batchNumber
			return
			
		}(idx)
		startIdx += adjustBatchSize // update StartIdx
		
	}

	return adjustBatchSize, nil
}

func (logic *Logic) SetNextBatchSize(predictTime float64, elapsedTime float64, previousBatchSize float64, currentBatchSize float64) (int, float64) {

    currentSpeed := currentBatchSize / elapsedTime

    nextBatch := int(currentSpeed * predictTime)

    if nextBatch < logic.config.MinBatchThreshold {
        nextBatch = logic.config.MinBatchThreshold
    } else if nextBatch > logic.config.MaxBatchSize {
        nextBatch = logic.config.MaxBatchSize 
    }

    nextPredictTime := (predictTime + elapsedTime) / 2

    return nextBatch, nextPredictTime
}


func (logic *Logic) Compute(device *types.Device, batchStartIdx int, batchEndIdx int) (float64, float64, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Compute",
	})

	logger.Debugf("received Compute()")

	type Response struct {
		ElapsedTime	float64	`json:"elapsed_time"`
		ComputeTime float64	`json:"compute_time"`
	}

	var response Response
	fullEndpoint := logic.GetFullEndpoint(device.IP, device.Port,device.Endpoint, batchStartIdx, batchEndIdx)
	
	client := resty.New()
	_, err := client.R().SetResult(&response).Get(fullEndpoint)
	
	if err != nil {
		fmt.Println("error in Compute():", err)
		return 0, 0, err
	}
	
	elapsedTime := response.ElapsedTime
	computeTime := response.ComputeTime
	return elapsedTime, computeTime, nil
}

// ywjang
func (logic *Logic) ScheduleJob(jobID string) error {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ScheduleJob",
	})
	
	logger.Debug("received ScheduleJob()")

	startTime := time.Now()
	fmt.Println(time.Now().Unix())
	job, err := logic.dbAdapter.GetJob(jobID)
	if err != nil {
		return err
	}
	
	// queue to hold available deviceID
	var deviceIDQueue Queue
	// map to hold previous and current latency
	deviceRecordMap := deviceRecord{}
	// enqueue deviceID into deviceIDQueue
	for _, deviceID := range job.DeviceIDList{
		deviceIDQueue.Enqueue(deviceID)
		deviceRecordMap[deviceID] = make([]float64, 7)
	}

	// initialize job, batch index
	jobStartIdx := job.StartIndex // start index for job; now set to zero for experiment, but may change later
	jobEndIdx := job.EndIndex // end index for job
	batchStartIdx := jobStartIdx // start index for batch, initialized as jobStartIdx
	batchEndIdx := 0 // end index for batch
	
	// job start index > job end index
	if jobStartIdx >= jobEndIdx {
		// job start index >= job end index -> terminate
		fmt.Println("schedule done: no need to work")
		return nil
	}

	// adjustBatchSize for entire time
	adjustBatchSize, err := logic.AdjustBatchSize(&deviceIDQueue, &deviceRecordMap, jobStartIdx, jobEndIdx)
		if err != nil {
			return err
	}
	fmt.Println("adjustBatchSize: ", adjustBatchSize)

	batchStartIdx = batchStartIdx + adjustBatchSize * len(job.DeviceIDList) // increase batch start index(adjustment stage completed)
	// channel for task completions or errors, change size if device number increases
	taskResults := make(chan error, 10) 
	
	iterNum := 0 // need to check if it's first batch of entire schedule

	prevBatchSize := float64(logic.config.MaxBatchSize)

	for {
		if deviceIDQueue.IsEmpty() {
			// If the queue is empty, wait for results from ongoing tasks before trying to dequeue again
			err := <-taskResults
			if err != nil {
				fmt.Println("Error processing batch:", err)
			}
			continue
		}

		// need to use channel instead of deviceQueue
		dID := deviceIDQueue.Dequeue()
		
		device, err := logic.dbAdapter.GetDevice(dID)
		if err != nil {
			deviceIDQueue.Enqueue(dID) // enqueue dID to try later
			return err
		}

		if (batchEndIdx == jobEndIdx) {
			fmt.Println("job done")
			break
		}

		batchSize := int(deviceRecordMap[dID][3]) // call batchSize from deviceRecord
		if iterNum > 0 {
			batchStartIdx = batchEndIdx
		}
		batchEndIdx = batchStartIdx + batchSize
		
		stopFlag := false 

		predictTime := deviceRecordMap[dID][1]
		if (batchStartIdx < jobEndIdx  && batchEndIdx > jobEndIdx){ 
			// batch start index is ok, but batch end index exceeds job end index -> too much batch size
			batchEndIdx = jobEndIdx // set batch end index to job end index
			batchSize = batchEndIdx - batchStartIdx
			predictTime = deviceRecordMap[dID][1] / deviceRecordMap[dID][2] * float64(batchSize)
			stopFlag = true
		}
		
		// batch size < lower threshold -> out of schedule
		if (batchSize <= logic.config.MinBatchThreshold) {
			fmt.Printf("device %s batchsize below %d, out of schedule!\n", dID, logic.config.MinBatchThreshold)
			continue
		}

		deviceRecordMap[dID][1] = predictTime
		deviceRecordMap[dID][2] = float64(batchSize) // update deviceRecordMap[2](current batch size) in the map
		deviceRecordMap[dID][4] += 1 // update batchNum
		deviceRecordMap[dID][5] = float64(batchStartIdx)
		deviceRecordMap[dID][6] = float64(batchEndIdx)

		// create goroutines
		if batchStartIdx < jobEndIdx {
			// send compute request to device
			// wg.Add(1)
			go func(dID string){
				// defer wg.Done()
				// batchsize lower than threshold
				// sleep for 10 secs -> return to job
				if (batchEndIdx - batchStartIdx < logic.config.TemporaryOutThreshold) { 
					fmt.Printf("device %s batchsize below %d, sleep for 10 secs\n", dID, logic.config.TemporaryOutThreshold)
					time.Sleep(10 * time.Second)
				}

				elapsedTime, _, err := logic.Compute(
					&device, 
					int(deviceRecordMap[dID][5]), 
					int(deviceRecordMap[dID][6]),
				)
				deviceRecordMap[dID][1] = elapsedTime // elapsedTime
				if err != nil {
					taskResults <- err
					return
				}

				// sbkwon 0326 수정
				nextBatchSize, predictTime := logic.SetNextBatchSize(
					deviceRecordMap[dID][0],
					deviceRecordMap[dID][1], 
					// prevTime,	
					prevBatchSize,
					float64(int(deviceRecordMap[dID][6])-int(deviceRecordMap[dID][5])),
				)
				fmt.Println("In Schedule Loop(deviceID, elapsedTime, batchSize, nextBatchSize, predictTime): ", dID, elapsedTime, int(int(deviceRecordMap[dID][6])-int(deviceRecordMap[dID][5])), nextBatchSize, predictTime, time.Now().Unix())
				deviceRecordMap[dID][0] = predictTime // predictTime
				deviceRecordMap[dID][2] = float64(batchSize) // current batch Size
				deviceRecordMap[dID][3] = float64(nextBatchSize) // next batch Size
				
				prevBatchSize = float64(batchSize) // sbkwon 0326
				// prevTime = elapsedTime // sbkwon 0326
								
				// job finished successfully
				if err == nil {
					taskResults <- nil
					deviceIDQueue.Enqueue(dID) // re-enqueue
				}
				
			}(dID)
			
			select {
			case err := <-taskResults:
				if err != nil {
					fmt.Println("Error processing batch:", err)
					deviceIDQueue.Enqueue(dID) // Possible re-enqueue on error
				}
			default:
				iterNum++
				// No waiting, continue to spawn new goroutines
			}
		}
		if stopFlag {
			break
		}
	}

	// update startIdx
	err = logic.dbAdapter.UpdateStartIndex(jobID, batchStartIdx)
	if err != nil {
        return err
    }

	// update Completed to true
    err = logic.dbAdapter.UpdateJobCompleted(jobID, true)
    if err != nil {
        return err
    }
	timeTaken := time.Since(startTime).Seconds()
	fmt.Println("time taken:", timeTaken)
	return nil
}
