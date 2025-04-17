package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/sensors"
)

func GetTemperatureData() (models.InfluxDbFields, error) {

	fields := make(models.InfluxDbFields)

	// Get and print temperatures
	temperatures, err := sensors.SensorsTemperatures()
	if err != nil {
		return nil, fmt.Errorf("Error getting temperatures: %v\n", err)
	}

	for _, temp := range temperatures {
		fields[temp.SensorKey] = temp.Temperature
	}

	return fields, nil
}
