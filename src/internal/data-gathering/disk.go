package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiscData() ([]models.InfluxDbField, error) {

	if diskUsage, err := disk.Usage("/"); err != nil {
		return nil, fmt.Errorf("Error getting disk usage: %v\n", err)
	} else {
		fields := []models.InfluxDbField{
			{Name: "total", Value: diskUsage.Total},
			{Name: "free", Value: diskUsage.Free},
			{Name: "used", Value: diskUsage.Used},
			{Name: "used_percent", Value: diskUsage.UsedPercent},
		}
		return fields, nil
	}
}
