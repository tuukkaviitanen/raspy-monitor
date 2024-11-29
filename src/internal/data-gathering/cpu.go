package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

func GetCpuData() (*models.CPUData, error) {
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

	parsedCPUData := &models.CPUData{
		LogicalCPUCount: logicalCount,
		TotalCPUUsage:   totalCPUUsage,
		CPUUsage:        cpuPercentages,
	}

	return parsedCPUData, nil
}
