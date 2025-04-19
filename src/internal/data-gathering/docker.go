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

func GetDockerData() (models.InfluxDbFields, error) {
	// Create docker stats command
	cmd := exec.Command("docker", "stats", "--format", "{{json .}}", "--no-stream")

	// Run command and wait for output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error creating Docker StdoutPipe:", err)
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

func handleStats(jsonStats string) models.InfluxDbFields {
	jsonLine := strings.Split(strings.TrimSpace(jsonStats), "\n")

	fields := models.InfluxDbFields{}

	for _, line := range jsonLine {
		var DockerStat models.DockerStat

		if err := json.Unmarshal([]byte(line), &DockerStat); err != nil {
			log.Printf("Error parsing JSON text '%s': %s\n", line, err)
			continue
		}

		parsedCPUPercentage, cpuParsingErr := parsePercentage(DockerStat.CPUPercentage)
		parsedMemPercentage, memParsingErr := parsePercentage(DockerStat.MemoryPercentage)
		parsedPidCount, pidParsingErr := strconv.Atoi(DockerStat.PIDs)

		if cpuParsingErr != nil {
			log.Printf("Error parsing CPU percentage: %s\n", cpuParsingErr)
			continue
		}

		if memParsingErr != nil {
			log.Printf("Error parsing Memory percentage: %s\n", memParsingErr)
			continue
		}

		if pidParsingErr != nil {
			log.Printf("Error parsing PID count: %s\n", pidParsingErr)
			continue
		}

		if _, exists := fields["cpu_usage_percentage"]; !exists {
			fields["cpu_usage_percentage"] = make([]models.InfluxDbTaggedValue, 0)
		}
		fields["cpu_usage_percentage"] = append(fields["cpu_usage_percentage"], models.InfluxDbTaggedValue{
			Value: parsedCPUPercentage,
			Tags: map[string]string{
				"container_name": DockerStat.Name,
			},
		})

		if _, exists := fields["memory_usage_percentage"]; !exists {
			fields["memory_usage_percentage"] = make([]models.InfluxDbTaggedValue, 0)
		}
		fields["memory_usage_percentage"] = append(fields["memory_usage_percentage"], models.InfluxDbTaggedValue{
			Value: parsedMemPercentage,
			Tags: map[string]string{
				"container_name": DockerStat.Name,
			},
		})

		if _, exists := fields["pid_count"]; !exists {
			fields["pid_count"] = make([]models.InfluxDbTaggedValue, 0)
		}
		fields["pid_count"] = append(fields["pid_count"], models.InfluxDbTaggedValue{
			Value: parsedPidCount,
			Tags: map[string]string{
				"container_name": DockerStat.Name,
			},
		})
	}

	return fields
}
