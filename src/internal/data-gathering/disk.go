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

	fields["total"] = []models.InfluxDbTaggedValue{{Value: diskUsage.Total}}
	fields["free"] = []models.InfluxDbTaggedValue{{Value: diskUsage.Free}}
	fields["used"] = []models.InfluxDbTaggedValue{{Value: diskUsage.Used}}
	fields["used_percent"] = []models.InfluxDbTaggedValue{{Value: diskUsage.UsedPercent}}

	return fields, nil
}
