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
