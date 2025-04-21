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

type InfluxDbTag struct {
	Name  string
	Value string
}

type InfluxDbField struct {
	Name  string
	Value any
	Tags  []InfluxDbTag
}

type InfluxDbMeasurement struct {
	Name   string
	Fields []InfluxDbField
}

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
