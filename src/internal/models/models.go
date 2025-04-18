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

type InfluxDbTaggedValue struct {
	Value any
	Tags  map[string]string
}

type InfluxDbFields map[string][]InfluxDbTaggedValue

type InfluxDbMeasurements map[string]InfluxDbFields

type DockerStat struct {
	Name             string `json:"Name"`
	ContainerId      string `json:"Container"`
	BlockIO          string `json:"BlockIO"`
	CPUPercentage    string `json:"CPUPerc"`
	MemoryPercentage string `json:"MemPerc"`
	MemoryUsage      string `json:"MemUsage"`
	NetIO            string `json:"NetIO"`
	PIDs             string `json:"PIDs"`
}
