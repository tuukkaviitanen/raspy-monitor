package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

func GetCpuData() (models.InfluxDbFields, error) {

	fields := make(models.InfluxDbFields)

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

	fields["total_cpu_usage"] = totalCPUUsage
	fields["logical_cpu_count"] = logicalCount

	for index, usage := range cpuPercentages {
		field := fmt.Sprintf("cpu_%d_usage", index)
		fields[field] = usage
	}

	return fields, nil
}
