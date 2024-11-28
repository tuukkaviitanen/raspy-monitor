package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

func main() {
	start := time.Now() // Start timing

	// Should only be fetched once on startup
	getHostInfo()

	getCpuInfo()

	getMemoryInfo()

	getDiscInfo()

	timestamp := time.Now()
	fmt.Printf("Timestamp: %s\n", timestamp.Format(time.RFC1123))

	elapsed := time.Since(start) // End timing
	fmt.Printf("Execution time: %s\n", elapsed)
}

func getHostInfo() {
	hostInfo, err := host.Info()
	if err != nil {
		fmt.Printf("Error getting host info: %v\n", err)
		return
	}

	fmt.Printf("Hostname: %s\n", hostInfo.Hostname)
	fmt.Printf("OS: %s\n", hostInfo.OS)
	fmt.Printf("Platform: %s\n", hostInfo.Platform)
	fmt.Printf("Platform Family: %s\n", hostInfo.PlatformFamily)
	fmt.Printf("Platform Version: %s\n", hostInfo.PlatformVersion)
	fmt.Printf("Kernel Version: %s\n", hostInfo.KernelVersion)
	fmt.Printf("Kernel Arch: %s\n", hostInfo.KernelArch)
	fmt.Printf("Boottime: %s\n", time.Unix(int64(hostInfo.BootTime), 0).Format(time.RFC1123))

}

func getCpuInfo() {
	logicalCount, err := cpu.Counts(true)
	if err != nil {
		fmt.Printf("Error getting CPU count: %v\n", err)
		return
	}

	fmt.Printf("Logical CPU Count: %d\n", logicalCount)
	cpuInterval := 1 * time.Second

	// Get CPU usage percentages for each CPU
	cpuPercentages, err := cpu.Percent(cpuInterval, true)
	if err != nil {
		fmt.Printf("Error getting CPU percentages: %v\n", err)
		return
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
}

func getMemoryInfo() {

	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, Cached: %v, Buffers: %v, Available: %v, Used: %v, UsedPercent:%f%%\n", v.Total, v.Free, v.Cached, v.Buffers, v.Available, v.Used, v.UsedPercent)
}

func getTemperatureInfo() {

	// Get and print temperatures
	temperatures, err := sensors.SensorsTemperatures()
	if err != nil {
		fmt.Printf("Error getting temperatures: %v\n", err)
		return
	}

	for _, temp := range temperatures {
		fmt.Printf("Sensor: %s, Temperature: %.2f°C\n", temp.SensorKey, temp.Temperature)
	}

	if len(temperatures) == 0 {
		fmt.Println("No temperature sensors found")
	}
}

func getDiscInfo() {

	diskUsage, err := disk.Usage("/")
	if err != nil {
		fmt.Printf("Error getting disk usage: %v\n", err)
		return
	}

	fmt.Printf("Disk Usage: Total: %v, Free: %v, Used: %v, UsedPercent: %.2f%%\n", diskUsage.Total, diskUsage.Free, diskUsage.Used, diskUsage.UsedPercent)

}
