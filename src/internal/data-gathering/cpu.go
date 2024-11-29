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

	cpuInterval := time.Second

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

	parsedCPUData := &models.CPUData{
		LogicalCPUCount: logicalCount,
		TotalCPUUsage:   totalCPUUsage,
		CPUUsage:        cpuPercentages,
	}

	return parsedCPUData, nil
}
