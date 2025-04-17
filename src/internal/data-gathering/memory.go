package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryData() (models.InfluxDbFields, error) {

	fields := make(models.InfluxDbFields)

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("Error getting memory info: %v\n", err)
	}

	fields["total"] = v.Total
	fields["free"] = v.Free
	fields["cached"] = v.Cached
	fields["buffers"] = v.Buffers
	fields["available"] = v.Available
	fields["used"] = v.Used
	fields["used_percent"] = v.UsedPercent

	return fields, nil
}
