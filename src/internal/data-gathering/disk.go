package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiscData() (*models.DiscData, error) {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("Error getting disk usage: %v\n", err)
	}

	fmt.Printf("Disk Usage: Total: %v, Free: %v, Used: %v, UsedPercent: %.2f%%\n", diskUsage.Total, diskUsage.Free, diskUsage.Used, diskUsage.UsedPercent)

	parsedDiscData := &models.DiscData{
		Total:       diskUsage.Total,
		Free:        diskUsage.Free,
		Used:        diskUsage.Used,
		UsedPercent: diskUsage.UsedPercent,
	}

	return parsedDiscData, nil
}
