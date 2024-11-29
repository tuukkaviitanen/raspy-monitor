package main

import (
	"fmt"
	datagathering "raspy-monitor/src/internal/data-gathering"
	"raspy-monitor/src/internal/influx"
	"time"
)

func main() {
	hostInfo, err := datagathering.GetHostData()
	if err != nil {
		fmt.Printf("Error getting host info: %v\n", err)
		return
	}

	fmt.Printf("Started monitoring hostname: %s\n", hostInfo.Hostname)
	fmt.Printf("OS: %s\n", hostInfo.OS)
	fmt.Printf("Platform: %s\n", hostInfo.Platform)
	fmt.Printf("Platform Family: %s\n", hostInfo.PlatformFamily)
	fmt.Printf("Platform Version: %s\n", hostInfo.PlatformVersion)
	fmt.Printf("Kernel Version: %s\n", hostInfo.KernelVersion)
	fmt.Printf("Kernel Arch: %s\n", hostInfo.KernelArch)
	fmt.Printf("Boot time: %s\n", hostInfo.BootTime.Format(time.RFC1123))

	// Run monitoring every second
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		go monitoringRun()
	}
}

func monitoringRun() {
	cpuData, err := datagathering.GetCpuData()
	if err != nil {
		fmt.Printf("Error getting CPU info: %v\n", err)
		return
	}

	memoryData, err := datagathering.GetMemoryData()
	if err != nil {
		fmt.Printf("Error getting memory info: %v\n", err)
		return
	}

	temperatureData, err := datagathering.GetTemperatureData()
	if err != nil {
		fmt.Printf("Error getting temperature info: %v\n", err)
		return
	}

	discData, err := datagathering.GetDiscData()
	if err != nil {
		fmt.Printf("Error getting disc info: %v\n", err)
		return
	}

	fmt.Println("Finished data collection")
	fmt.Println("Writing data to InfluxDB...")

	influx.WriteSystemDataToInflux(cpuData, memoryData, temperatureData, discData)

	fmt.Println("Data writing to influxDB finished")
}
