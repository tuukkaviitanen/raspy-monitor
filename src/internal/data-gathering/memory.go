package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryData() (*models.MemoryData, error) {

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("Error getting memory info: %v\n", err)
	}

	parsedMemoryData := &models.MemoryData{
		Total:       v.Total,
		Free:        v.Free,
		Cached:      v.Cached,
		Buffers:     v.Buffers,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}

	return parsedMemoryData, nil
}
