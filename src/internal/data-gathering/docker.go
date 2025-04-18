package datagathering

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"
	"raspy-monitor/src/internal/models"
	"regexp"
	"strconv"
	"strings"
)

func StreamDockerData(callback func(models.InfluxDbMeasurements)) {
	// Create docker stats command
	cmd := exec.Command("docker", "stats", "--format", "{{json .}}")

	// Get the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating Docker StdoutPipe:", err)
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting Docker stats command:", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// Launching a goroutine to handle each line of output concurrently
		go handleStat(scanner.Text(), callback)
	}
}

// Define a regular expression to match ANSI escape codes
const ansiEscapePattern = `\x1b\[[0-9;]*[a-zA-Z]`

var ansiEscapeRegex = regexp.MustCompile(ansiEscapePattern)

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

func handleStat(jsonStats string, callback func(models.InfluxDbMeasurements)) {

	// Remove ANSI escape codes from the JSON string
	cleanedJsonString := ansiEscapeRegex.ReplaceAllString(jsonStats, "")

	// Define a variable to hold the parsed JSON DockerStat
	var DockerStat models.DockerStat

	// Parse the JSON string into the map
	if err := json.Unmarshal([]byte(cleanedJsonString), &DockerStat); err != nil {
		log.Printf("Error parsing JSON text '%s': %s\n", cleanedJsonString, err)
		return
	}

	parsedCPUPercentage, cpuParsingErr := parsePercentage(DockerStat.CPUPercentage)
	parsedMemPercentage, memParsingErr := parsePercentage(DockerStat.MemoryPercentage)
	parsedPidCount, pidParsingErr := strconv.Atoi(DockerStat.PIDs)

	if cpuParsingErr != nil {
		log.Printf("Error parsing CPU percentage: %s\n", cpuParsingErr)
		return
	}

	if memParsingErr != nil {
		log.Printf("Error parsing Memory percentage: %s\n", memParsingErr)
		return
	}

	if pidParsingErr != nil {
		log.Printf("Error parsing PID count: %s\n", pidParsingErr)
		return
	}

	// Create a map to hold the measurements
	measurements := models.InfluxDbMeasurements{
		"docker_data": {
			"cpu_usage_percentage": {
				models.InfluxDbTaggedValue{
					Value: parsedCPUPercentage,
					Tags: map[string]string{
						"container_name": DockerStat.Name,
					},
				},
			},
			"memory_usage_percentage": {
				models.InfluxDbTaggedValue{
					Value: parsedMemPercentage,
					Tags: map[string]string{
						"container_name": DockerStat.Name,
					},
				},
			},
			"pid_count": {
				models.InfluxDbTaggedValue{
					Value: parsedPidCount,
					Tags: map[string]string{
						"container_name": DockerStat.Name,
					},
				},
			},
		},
	}
	callback(measurements)
}
