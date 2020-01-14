package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type serviceDiscovery struct {
	once *sync.Once
}

var ServiceDiscovery = newServiceDiscovery()

func newServiceDiscovery() *serviceDiscovery {
	return &serviceDiscovery{
		once: new(sync.Once),
	}
}

func (s *serviceDiscovery) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	s.once.Do(func() {
		public.Path("/health").HandlerFunc(s.healthCheck).Methods("GET")
		public.Path("/disk").HandlerFunc(s.diskCheck).Methods("GET")
		public.Path("/cpu").HandlerFunc(s.cpuCheck).Methods("GET")
		public.Path("/ram").HandlerFunc(s.ramCheck).Methods("GET")
	})
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// HealthCheck shows `OK` as the ping-pong result.
func (s *serviceDiscovery) healthCheck(w http.ResponseWriter, r *http.Request) {
	message := "OK"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("\n" + message))
}

// DiskCheck checks the disk usage.
func (s *serviceDiscovery) diskCheck(w http.ResponseWriter, r *http.Request) {
	u, _ := disk.Usage("/")

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusOK
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	w.WriteHeader(status)
	w.Write([]byte("\n" + message))
}

// CPUCheck checks the cpu usage.
func (s *serviceDiscovery) cpuCheck(w http.ResponseWriter, r *http.Request) {
	cores, _ := cpu.Counts(false)

	a, _ := load.Avg()
	l1 := a.Load1
	l5 := a.Load5
	l15 := a.Load15

	status := http.StatusOK
	text := "OK"

	if l5 >= float64(cores-1) {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)
	w.WriteHeader(status)
	w.Write([]byte("\n" + message))
}

// RAMCheck checks the disk usage.
func (s *serviceDiscovery) ramCheck(w http.ResponseWriter, r *http.Request) {
	u, _ := mem.VirtualMemory()

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	w.WriteHeader(status)
	w.Write([]byte("\n" + message))
}
