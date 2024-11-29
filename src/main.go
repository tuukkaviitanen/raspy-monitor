package main

import (
	"fmt"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

var (
	INFLUXDB_URL    = os.Getenv("INFLUXDB_URL")
	INFLUXDB_TOKEN  = os.Getenv("INFLUXDB_TOKEN")
	INFLUXDB_ORG    = os.Getenv("INFLUXDB_ORG")
	INFLUXDB_BUCKET = os.Getenv("INFLUXDB_BUCKET")
)

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
	SensorTemperatures map[string]float64
}

type DiscData struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

func main() {
	hostInfo, err := getHostData()
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
	cpuData, err := getCpuData()
	if err != nil {
		fmt.Printf("Error getting CPU info: %v\n", err)
		return
	}

	memoryData, err := getMemoryData()
	if err != nil {
		fmt.Printf("Error getting memory info: %v\n", err)
		return
	}

	temperatureData, err := getTemperatureData()
	if err != nil {
		fmt.Printf("Error getting temperature info: %v\n", err)
		return
	}

	discData, err := getDiscData()
	if err != nil {
		fmt.Printf("Error getting disc info: %v\n", err)
		return
	}

	fmt.Println("Finished data collection")
	fmt.Println("Writing data to InfluxDB...")

	saveDataToInflux(cpuData, memoryData, temperatureData, discData)

	fmt.Println("Data writing to influxDB finished")

}

func getHostData() (*HostInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		fmt.Printf("Error getting host info: %v\n", err)
		return nil, err
	}

	parsedHostInfo := &HostInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformFamily:  hostInfo.PlatformFamily,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		KernelArch:      hostInfo.KernelArch,
		BootTime:        time.Unix(int64(hostInfo.BootTime), 0),
	}

	return parsedHostInfo, nil
}

func getCpuData() (*CPUData, error) {
	logicalCount, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("Error getting CPU count: %v\n", err)
	}

	fmt.Printf("Logical CPU Count: %d\n", logicalCount)
	cpuInterval := 1 * time.Second

	// Get CPU usage percentages for each CPU
	cpuPercentages, err := cpu.Percent(cpuInterval, true)
	if err != nil {
		return nil, fmt.Errorf("Error getting CPU percentages: %v\n", err)
	}

	// Print total CPU usage (average of all CPUs)
	totalCPUUsage := 0.0
	for _, percentage := range cpuPercentages {
		totalCPUUsage += percentage
	}
	totalCPUUsage /= float64(len(cpuPercentages))
	fmt.Printf("Total CPU Usage: %.2f%%\n", totalCPUUsage)

	// Print CPU usage for each CPU
	for i, percentage := range cpuPercentages {
		fmt.Printf("CPU %d: %.2f%%\n", i, percentage)
	}

	parsedCPUData := &CPUData{
		LogicalCPUCount: logicalCount,
		TotalCPUUsage:   totalCPUUsage,
		CPUUsage:        cpuPercentages,
	}

	return parsedCPUData, nil
}

func getMemoryData() (*MemoryData, error) {

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("Error getting memory info: %v\n", err)
	}

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, Cached: %v, Buffers: %v, Available: %v, Used: %v, UsedPercent:%f%%\n", v.Total, v.Free, v.Cached, v.Buffers, v.Available, v.Used, v.UsedPercent)

	parsedMemoryData := &MemoryData{
		Total:       v.Total,
		Free:        v.Free,
		Cached:      v.Cached,
		Buffers:     v.Buffers,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}

	return parsedMemoryData, nil
}

func getTemperatureData() (*TemperatureData, error) {

	// Get and print temperatures
	temperatures, err := sensors.SensorsTemperatures()
	if err != nil {
		return nil, fmt.Errorf("Error getting temperatures: %v\n", err)
	}

	for _, temp := range temperatures {
		fmt.Printf("Sensor: %s, Temperature: %.2f°C\n", temp.SensorKey, temp.Temperature)
	}

	if len(temperatures) == 0 {
		fmt.Println("No temperature sensors found")
	}

	sensorMap := make(map[string]float64)

	for _, temp := range temperatures {
		sensorMap[temp.SensorKey] = temp.Temperature
	}

	parsedTemperatureData := &TemperatureData{
		SensorTemperatures: sensorMap,
	}

	return parsedTemperatureData, nil
}

func getDiscData() (*DiscData, error) {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("Error getting disk usage: %v\n", err)
	}

	fmt.Printf("Disk Usage: Total: %v, Free: %v, Used: %v, UsedPercent: %.2f%%\n", diskUsage.Total, diskUsage.Free, diskUsage.Used, diskUsage.UsedPercent)

	parsedDiscData := &DiscData{
		Total:       diskUsage.Total,
		Free:        diskUsage.Free,
		Used:        diskUsage.Used,
		UsedPercent: diskUsage.UsedPercent,
	}

	return parsedDiscData, nil
}

func saveDataToInflux(cpuData *CPUData, memoryData *MemoryData, temperatureData *TemperatureData, discData *DiscData) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)

	writeAPI := client.WriteAPI(INFLUXDB_ORG, INFLUXDB_BUCKET)

	// Start a goroutine to handle errors
	go func() {
		for err := range writeAPI.Errors() {
			fmt.Printf("Write error: %v\n", err)
		}
	}()

	timestamp := time.Now()

	if cpuData != nil {
		point := influxdb2.NewPointWithMeasurement("cpu_data").
			AddField("total_cpu_usage", cpuData.TotalCPUUsage).
			SetTime(timestamp)

		for index, usage := range cpuData.CPUUsage {
			point.AddField(fmt.Sprintf("cpu_%d_usage", index), usage)
		}

		writeAPI.WritePoint(point)
	}

	if memoryData != nil {
		point := influxdb2.NewPointWithMeasurement("memory_data").
			AddField("total", memoryData.Total).
			AddField("free", memoryData.Free).
			AddField("cached", memoryData.Cached).
			AddField("buffers", memoryData.Buffers).
			AddField("available", memoryData.Available).
			AddField("used", memoryData.Used).
			AddField("used_percent", memoryData.UsedPercent).
			SetTime(timestamp)

		writeAPI.WritePoint(point)
	}

	if temperatureData != nil && len(temperatureData.SensorTemperatures) > 0 {
		point := influxdb2.NewPointWithMeasurement("temperature_data").
			SetTime(timestamp)

		for sensor, temperature := range temperatureData.SensorTemperatures {
			point.AddField(sensor, temperature)
		}

		writeAPI.WritePoint(point)
	}

	if discData != nil {
		point := influxdb2.NewPointWithMeasurement("disc_data").
			AddField("total", discData.Total).
			AddField("free", discData.Free).
			AddField("used", discData.Used).
			AddField("used_percent", discData.UsedPercent).
			SetTime(timestamp)

		writeAPI.WritePoint(point)
	}

	writeAPI.Flush()

	defer client.Close()
}
