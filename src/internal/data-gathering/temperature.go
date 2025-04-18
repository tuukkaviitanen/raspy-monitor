package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/sensors"
)

func GetTemperatureData() (models.InfluxDbFields, error) {

	fields := models.InfluxDbFields{}

	// Get and print temperatures
	temperatures, err := sensors.SensorsTemperatures()
	if err != nil {
		return nil, fmt.Errorf("Error getting temperatures: %v\n", err)
	}

	temperatureField := []models.InfluxDbTaggedValue{}

	for _, temp := range temperatures {
		temperatureField = append(temperatureField, models.InfluxDbTaggedValue{Value: temp.Temperature, Tags: map[string]string{"sensor": temp.SensorKey}})
	}

	fields["temperature"] = temperatureField

	return fields, nil
}
