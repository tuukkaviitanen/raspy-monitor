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

	parsedDiscData := &models.DiscData{
		Total:       diskUsage.Total,
		Free:        diskUsage.Free,
		Used:        diskUsage.Used,
		UsedPercent: diskUsage.UsedPercent,
	}

	return parsedDiscData, nil
}
