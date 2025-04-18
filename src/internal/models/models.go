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

type InfluxDbValue any

type InfluxDbFields map[string]InfluxDbValue

type InfluxDbMeasurements map[string]InfluxDbFields
