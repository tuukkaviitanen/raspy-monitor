package main

import (
	"log"
	datagathering "raspy-monitor/src/internal/data-gathering"
	"raspy-monitor/src/internal/influx"
	"raspy-monitor/src/internal/models"
	"time"
)

func main() {
	{
		hostInfo, err := datagathering.GetHostData()
		if err != nil {
			log.Fatalf("Error getting host info: %v\n", err)
			return
		}

		log.Printf("Started monitoring hostname %s at %s\n", hostInfo.Hostname, time.Now().Format(time.RFC1123))
		log.Printf("OS: %s\n", hostInfo.OS)
		log.Printf("Platform: %s\n", hostInfo.Platform)
		log.Printf("Platform Family: %s\n", hostInfo.PlatformFamily)
		log.Printf("Platform Version: %s\n", hostInfo.PlatformVersion)
		log.Printf("Kernel Version: %s\n", hostInfo.KernelVersion)
		log.Printf("Kernel Arch: %s\n", hostInfo.KernelArch)
		log.Printf("Boot time: %s\n", hostInfo.BootTime.Format(time.RFC1123))
	} // No need to keep hostInfo in memory so closing the scope

	// Run monitoring every second
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		go monitoringRun()
	}
}

func monitoringRun() {
	log.Println("Starting data collection...")

	measurements := make(models.InfluxDbMeasurements)

	cpuData, err := datagathering.GetCpuData()
	if err != nil {
		log.Printf("Error getting CPU info: %v\n", err)
	} else {
		measurements["cpu_data"] = cpuData
	}

	memoryData, err := datagathering.GetMemoryData()
	if err != nil {
		log.Printf("Error getting memory info: %v\n", err)
	} else {
		measurements["memory_data"] = memoryData
	}

	temperatureData, err := datagathering.GetTemperatureData()
	if err != nil {
		log.Printf("Error getting temperature info: %v\n", err)
	} else {
		measurements["temperature_data"] = temperatureData
	}

	discData, err := datagathering.GetDiscData()
	if err != nil {
		log.Printf("Error getting disc info: %v\n", err)
	} else {
		measurements["disc_data"] = discData
	}

	log.Println("Finished data collection")

	log.Println("Writing data to InfluxDB...")

	influx.WriteSystemDataToInflux(measurements)

	log.Println("Writing data to influxDB done")
}
