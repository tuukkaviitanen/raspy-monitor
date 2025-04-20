package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryData() ([]models.InfluxDbField, error) {
	if v, err := mem.VirtualMemory(); err != nil {
		return nil, fmt.Errorf("Error getting memory info: %v\n", err)
	} else {
		fields := []models.InfluxDbField{
			{Name: "total", Value: v.Total},
			{Name: "free", Value: v.Free},
			{Name: "cached", Value: v.Cached},
			{Name: "buffers", Value: v.Buffers},
			{Name: "available", Value: v.Available},
			{Name: "used", Value: v.Used},
			{Name: "used_percent", Value: v.UsedPercent},
		}
		return fields, nil
	}
}
