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

	for _, temp := range temperatures {
		fields["temperature"] = []models.InfluxDbTaggedValue{{Value: temp.Temperature, Tags: map[string]string{"sensor": temp.SensorKey}}}
	}

	return fields, nil
}
