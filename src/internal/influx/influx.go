package influx

import (
	"fmt"
	"log"
	"os"
	"raspy-monitor/src/internal/models"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	INFLUXDB_URL    = os.Getenv("INFLUXDB_URL")
	INFLUXDB_TOKEN  = os.Getenv("INFLUXDB_TOKEN")
	INFLUXDB_ORG    = os.Getenv("INFLUXDB_ORG")
	INFLUXDB_BUCKET = os.Getenv("INFLUXDB_BUCKET")
)

func WriteSystemDataToInflux(cpuData *models.CPUData, memoryData *models.MemoryData, temperatureData *models.TemperatureData, discData *models.DiscData) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)

	writeAPI := client.WriteAPI(INFLUXDB_ORG, INFLUXDB_BUCKET)

	// Start a goroutine to handle errors
	go func() {
		for err := range writeAPI.Errors() {
			log.Printf("InfluxDB write error: %v\n", err)
		}
	}()

	timestamp := time.Now()

	if cpuData != nil {
		point := influxdb2.NewPointWithMeasurement("cpu_data").
			AddField("total_cpu_usage", cpuData.TotalCPUUsage).
			SetTime(timestamp)

		for index, usage := range cpuData.CPUUsage {
			point.AddField(fmt.Sprintf("cpu_%d_usage", index), usage)
		}

		writeAPI.WritePoint(point)
	}

	if memoryData != nil {
		point := influxdb2.NewPointWithMeasurement("memory_data").
			AddField("total", memoryData.Total).
			AddField("free", memoryData.Free).
			AddField("cached", memoryData.Cached).
			AddField("buffers", memoryData.Buffers).
			AddField("available", memoryData.Available).
			AddField("used", memoryData.Used).
			AddField("used_percent", memoryData.UsedPercent).
			SetTime(timestamp)

		writeAPI.WritePoint(point)
	}

	if temperatureData != nil && len(temperatureData.Temperatures) > 0 {
		point := influxdb2.NewPointWithMeasurement("temperature_data").
			SetTime(timestamp)

		for sensor, temperature := range temperatureData.Temperatures {
			point.AddField(sensor, temperature)
		}

		writeAPI.WritePoint(point)
	}

	if discData != nil {
		point := influxdb2.NewPointWithMeasurement("disc_data").
			AddField("total", discData.Total).
			AddField("free", discData.Free).
			AddField("used", discData.Used).
			AddField("used_percent", discData.UsedPercent).
			SetTime(timestamp)

		writeAPI.WritePoint(point)
	}

	writeAPI.Flush()

	defer client.Close()
}
