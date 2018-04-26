package apihelper

import (
	"os"
	"strings"
	"time"

	"github.com/json-iterator/go"
	"gitlab.com/project-d-collab/dhelpers"
)

// WorkerJobInformation contains information about one Job at a owkrer
type WorkerJobInformation struct {
	Function string
	Next     time.Time
	Prev     time.Time
}

// WorkerStatus contains information about a worker
type WorkerStatus struct {
	Available bool
	Entries   []WorkerJobInformation
}

// GetWorkerStatus returns information about all workers
// the addresses are read from WORKER_ADDRESSES, split using commas
func GetWorkerStatus() (status map[string]WorkerStatus) {
	status = make(map[string]WorkerStatus)
	workerAddresses := os.Getenv("WORKER_ADDRESSES")
	for _, workerAddress := range strings.Split(workerAddresses, ",") {
		workerAddress = strings.TrimSpace(workerAddress)
		data, err := dhelpers.NetGet(workerAddress + "/stats/cron")
		if err != nil {
			status[workerAddress] = WorkerStatus{
				Available: false,
			}
			continue
		}
		var entries []WorkerJobInformation
		err = jsoniter.Unmarshal(data, &entries)
		dhelpers.CheckErr(err)
		status[workerAddress] = WorkerStatus{
			Available: true,
			Entries:   entries,
		}
	}
	return status
}
