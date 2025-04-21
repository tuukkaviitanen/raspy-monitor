package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"

	"github.com/shirou/gopsutil/v4/sensors"
)

func GetTemperatureData() ([]models.InfluxDbField, error) {

	// Get and print temperatures
	if temperatures, err := sensors.SensorsTemperatures(); err != nil {
		return nil, fmt.Errorf("Error getting temperatures: %v\n", err)
	} else {

		fields := []models.InfluxDbField{}

		for _, temp := range temperatures {
			fields = append(fields, models.InfluxDbField{
				Name:  "temperature",
				Value: temp.Temperature,
				Tags:  []models.InfluxDbTag{{Name: "sensor", Value: temp.SensorKey}},
			})
		}

		return fields, nil
	}
}
