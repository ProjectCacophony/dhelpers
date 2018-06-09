package apihelper

import (
	"os"
	"strings"
	"time"

	"runtime"

	"net/http"

	"github.com/json-iterator/go"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/metrics"
	"gitlab.com/Cacophony/dhelpers/net"
)

// WorkerJobInformation contains information about one Job at a Worker
type WorkerJobInformation struct {
	Function string
	Next     time.Time
	Prev     time.Time
}

// GatewayEventInformation contains information about the events received by a Gateway
type GatewayEventInformation struct {
	EventsDiscarded                int64
	EventsGuildCreate              int64
	EventsGuildUpdate              int64
	EventsGuildDelete              int64
	EventsGuildMemberAdd           int64
	EventsGuildMemberUpdate        int64
	EventsGuildMemberRemove        int64
	EventsGuildMembersChunk        int64
	EventsGuildRoleCreate          int64
	EventsGuildRoleUpdate          int64
	EventsGuildRoleDelete          int64
	EventsGuildEmojisUpdate        int64
	EventsChannelCreate            int64
	EventsChannelUpdate            int64
	EventsChannelDelete            int64
	EventsMessageCreate            int64
	EventsMessageUpdate            int64
	EventsMessageDelete            int64
	EventsPresenceUpdate           int64
	EventsChannelPinsUpdate        int64
	EventsGuildBanAdd              int64
	EventsGuildBanRemove           int64
	EventsMessageReactionAdd       int64
	EventsMessageReactionRemove    int64
	EventsMessageReactionRemoveAll int64
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

// WorkerStatus contains information about a Worker
type WorkerStatus struct {
	Available bool
	Entries   []WorkerJobInformation
	Service   ServiceInformation
}

// Render renders the WorkerStatus for a network response, required to satisfy Chi interface
func (s WorkerStatus) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GatewayStatus contains information about a Gateway
type GatewayStatus struct {
	Available bool
	Service   ServiceInformation
	Events    GatewayEventInformation
}

// Render renders the WorkerStatus for a network response, required to satisfy Chi interface
func (s GatewayStatus) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// SqsProcessorStatus contains information about a SqsProcessor
type SqsProcessorStatus struct {
	Available bool
	Service   ServiceInformation
}

// Render renders the WorkerStatus for a network response, required to satisfy Chi interface
func (s SqsProcessorStatus) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
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
		data, err := net.Get(workerAddress + "/stats")
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
		data, err := net.Get(gatewayAddress + "/stats")
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

// ReadSqsProcessorStatus returns information about all workers
// the addresses are read from SQSPROCESSOR_ADDRESSES, split using commas
func ReadSqsProcessorStatus() (stats map[string]SqsProcessorStatus) {
	stats = make(map[string]SqsProcessorStatus)
	sqsProcessorAddresses := os.Getenv("SQSPROCESSOR_ADDRESSES")
	for _, sqsProcessorAddress := range strings.Split(sqsProcessorAddresses, ",") {
		sqsProcessorAddress = strings.TrimSpace(sqsProcessorAddress)
		data, err := net.Get(sqsProcessorAddress + "/stats")
		if err != nil {
			stats[sqsProcessorAddress] = SqsProcessorStatus{
				Available: false,
			}
			continue
		}
		var status SqsProcessorStatus
		err = jsoniter.Unmarshal(data, &status)
		dhelpers.CheckErr(err)
		stats[sqsProcessorAddress] = status
	}
	return stats
}
