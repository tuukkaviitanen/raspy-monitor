package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryData() (models.InfluxDbFields, error) {

	fields := models.InfluxDbFields{}

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("Error getting memory info: %v\n", err)
	}

	fields["total"] = []models.InfluxDbTaggedValue{{Value: v.Total}}
	fields["free"] = []models.InfluxDbTaggedValue{{Value: v.Free}}
	fields["cached"] = []models.InfluxDbTaggedValue{{Value: v.Cached}}
	fields["buffers"] = []models.InfluxDbTaggedValue{{Value: v.Buffers}}
	fields["available"] = []models.InfluxDbTaggedValue{{Value: v.Available}}
	fields["used"] = []models.InfluxDbTaggedValue{{Value: v.Used}}
	fields["used_percent"] = []models.InfluxDbTaggedValue{{Value: v.UsedPercent}}

	return fields, nil
}
