package datagathering

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"raspy-monitor/src/internal/models"
	"strconv"
	"strings"
)

func GetDockerData() ([]models.InfluxDbField, error) {
	// Create docker stats command
	cmd := exec.Command("docker", "stats", "--format", "{{json .}}", "--no-stream")

	// Run command and wait for output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error executing docker stats: %w", err)
	}

	outputStr := string(output)

	fields := handleStats(outputStr)

	return fields, nil
}

func parsePercentage(percentageStr string) (float64, error) {
	// Remove the '%' sign from the end of the string
	cleanedStr := strings.TrimSuffix(percentageStr, "%")

	// Parse the cleaned string to float64
	parsedFloat, err := strconv.ParseFloat(cleanedStr, 64)
	if err != nil {
		return 0, err
	}

	return parsedFloat, nil
}

func handleStats(jsonStats string) []models.InfluxDbField {
	jsonLine := strings.Split(strings.TrimSpace(jsonStats), "\n")

	fields := []models.InfluxDbField{}

	for _, line := range jsonLine {
		var DockerStat models.DockerStat

		if err := json.Unmarshal([]byte(line), &DockerStat); err != nil {
			log.Printf("Error parsing JSON text '%s': %s\n", line, err)
			continue
		}

		if parsedCPUPercentage, err := parsePercentage(DockerStat.CPUPercentage); err != nil {
			log.Printf("Error parsing CPU percentage: %s\n", err)
			continue
		} else {
			fields = append(fields, models.InfluxDbField{
				Name:  "cpu_usage_percentage",
				Value: parsedCPUPercentage,
				Tags: []models.InfluxDbTag{{
					Name: "container_name", Value: DockerStat.Name,
				}},
			})
		}

		if parsedMemPercentage, err := parsePercentage(DockerStat.MemoryPercentage); err != nil {
			log.Printf("Error parsing Memory percentage: %s\n", err)
			continue
		} else {
			fields = append(fields, models.InfluxDbField{
				Name:  "memory_usage_percentage",
				Value: parsedMemPercentage,
				Tags: []models.InfluxDbTag{{
					Name: "container_name", Value: DockerStat.Name,
				}},
			})
		}

		if parsedPidCount, err := strconv.Atoi(DockerStat.PIDs); err != nil {
			log.Printf("Error parsing Memory percentage: %s\n", err)
			continue
		} else {
			fields = append(fields, models.InfluxDbField{
				Name:  "pid_count",
				Value: parsedPidCount,
				Tags: []models.InfluxDbTag{{
					Name: "container_name", Value: DockerStat.Name,
				}},
			})
		}

	}

	return fields
}
