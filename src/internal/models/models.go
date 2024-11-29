package models

import "time"

type HostInfo struct {
	Hostname        string
	OS              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
	BootTime        time.Time
}

type CPUData struct {
	LogicalCPUCount int
	TotalCPUUsage   float64
	CPUUsage        []float64
}

type MemoryData struct {
	Total       uint64
	Free        uint64
	Cached      uint64
	Buffers     uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
}

type TemperatureData struct {
	Temperatures map[string]float64
}

type DiscData struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}
