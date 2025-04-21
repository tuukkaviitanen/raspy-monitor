package main

import (
	"log"
	datagathering "raspy-monitor/src/internal/data-gathering"
	"raspy-monitor/src/internal/influx"
	"raspy-monitor/src/internal/models"
	"sync"
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

	measurementChannel := make(chan models.InfluxDbMeasurement)
	var waitGroup sync.WaitGroup

	gatherDataAsynchronously("cpu_data", measurementChannel, &waitGroup, datagathering.GetCpuData)
	gatherDataAsynchronously("memory_data", measurementChannel, &waitGroup, datagathering.GetMemoryData)
	gatherDataAsynchronously("temperature_data", measurementChannel, &waitGroup, datagathering.GetTemperatureData)
	gatherDataAsynchronously("disc_data", measurementChannel, &waitGroup, datagathering.GetDiscData)
	gatherDataAsynchronously("docker_data", measurementChannel, &waitGroup, datagathering.GetDockerData)

	// Close the channel asynchronously once all the data is gathered
	go func() {
		waitGroup.Wait()
		close(measurementChannel)
	}()

	measurements := []models.InfluxDbMeasurement{}

	// Waits and collects gathered data from the channel until it is closed
	for measurement := range measurementChannel {
		measurements = append(measurements, measurement)
	}

	log.Println("Finished data collection, writing data to InfluxDB...")

	influx.WriteSystemDataToInflux(measurements)

	log.Println("Data written to InfluxDB successfully")
}

func gatherDataAsynchronously(measurementName string, measurementChanel chan<- models.InfluxDbMeasurement, wg *sync.WaitGroup, gatherFunc func() ([]models.InfluxDbField, error)) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		if gatheredData, err := gatherFunc(); err != nil {
			log.Printf("Error gathering %s: %v\n", measurementName, err)
		} else {
			measurementChanel <- models.InfluxDbMeasurement{Name: measurementName, Fields: gatheredData}
		}
	}()
}
