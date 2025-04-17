package influx

import (
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

func WriteSystemDataToInflux(measurements models.InfluxDbMeasurements) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)

	writeAPI := client.WriteAPI(INFLUXDB_ORG, INFLUXDB_BUCKET)

	// Start a goroutine to handle errors
	go func() {
		for err := range writeAPI.Errors() {
			log.Printf("InfluxDB write error: %v\n", err)
		}
	}()

	timestamp := time.Now()

	for measurement, fields := range measurements {
		if len(fields) == 0 {
			log.Printf("No fields to write for measurement %s\n", measurement)
			continue
		}

		point := influxdb2.NewPointWithMeasurement(measurement).SetTime(timestamp)

		for field, value := range fields {
			point.AddField(field, value)
		}
		writeAPI.WritePoint(point)
	}

	writeAPI.Flush()

	defer client.Close()
}
