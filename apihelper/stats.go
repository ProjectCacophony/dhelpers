package apihelper

import (
	"os"
	"strings"
	"time"

	"runtime"

	"github.com/json-iterator/go"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/metrics"
)

// WorkerJobInformation contains information about one Job at a owkrer
type WorkerJobInformation struct {
	Function string
	Next     time.Time
	Prev     time.Time
}

// ServiceInformation contains general information about a service
type ServiceInformation struct {
	Heap       uint64
	Sys        uint64
	Coroutines int
	GC         uint64
	Launch     time.Time
	Go         string
}

// WorkerStatus contains information about a worker
type WorkerStatus struct {
	Available bool
	Entries   []WorkerJobInformation
	Service   ServiceInformation
}

// GatewayStatus contains information about a gateway
type GatewayStatus struct {
	Available bool
	Service   ServiceInformation
}

// GenerateServiceInformation generates general information about a go processor
func GenerateServiceInformation() (information ServiceInformation) {
	var ram runtime.MemStats
	runtime.ReadMemStats(&ram)
	information.Heap = ram.Alloc
	information.Sys = ram.Sys
	information.Coroutines = runtime.NumGoroutine()
	information.GC = ram.TotalAlloc
	information.Launch = time.Unix(metrics.Uptime.Value(), 0)
	information.Go = strings.Replace(runtime.Version(), "go", "", 1)
	return information
}

// ReadWorkerStatus returns information about all workers
// the addresses are read from WORKER_ADDRESSES, split using commas
func ReadWorkerStatus() (stats map[string]WorkerStatus) {
	stats = make(map[string]WorkerStatus)
	workerAddresses := os.Getenv("WORKER_ADDRESSES")
	for _, workerAddress := range strings.Split(workerAddresses, ",") {
		workerAddress = strings.TrimSpace(workerAddress)
		data, err := dhelpers.NetGet(workerAddress + "/stats")
		if err != nil {
			stats[workerAddress] = WorkerStatus{
				Available: false,
			}
			continue
		}
		var status WorkerStatus
		err = jsoniter.Unmarshal(data, &status)
		dhelpers.CheckErr(err)
		stats[workerAddress] = status
	}
	return stats
}

// ReadGatewayStatus returns information about all workers
// the addresses are read from WORKER_ADDRESSES, split using commas
func ReadGatewayStatus() (stats map[string]GatewayStatus) {
	stats = make(map[string]GatewayStatus)
	gatewayAddresses := os.Getenv("GATEWAY_ADDRESSES")
	for _, gatewayAddress := range strings.Split(gatewayAddresses, ",") {
		gatewayAddress = strings.TrimSpace(gatewayAddress)
		data, err := dhelpers.NetGet(gatewayAddress + "/stats")
		if err != nil {
			stats[gatewayAddress] = GatewayStatus{
				Available: false,
			}
			continue
		}
		var status GatewayStatus
		err = jsoniter.Unmarshal(data, &status)
		dhelpers.CheckErr(err)
		stats[gatewayAddress] = status
	}
	return stats
}
