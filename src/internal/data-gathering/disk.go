package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiscData() (models.InfluxDbFields, error) {
	fields := models.InfluxDbFields{}

	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("Error getting disk usage: %v\n", err)
	}

	fields["total"] = diskUsage.Total
	fields["free"] = diskUsage.Free
	fields["used"] = diskUsage.Used
	fields["used_percent"] = diskUsage.UsedPercent

	return fields, nil
}
