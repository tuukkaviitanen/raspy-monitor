package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/sensors"
)

func GetTemperatureData() (*models.TemperatureData, error) {

	// Get and print temperatures
	temperatures, err := sensors.SensorsTemperatures()
	if err != nil {
		return nil, fmt.Errorf("Error getting temperatures: %v\n", err)
	}

	sensorMap := make(map[string]float64)

	for _, temp := range temperatures {
		sensorMap[temp.SensorKey] = temp.Temperature
	}

	parsedTemperatureData := &models.TemperatureData{
		Temperatures: sensorMap,
	}

	return parsedTemperatureData, nil
}
