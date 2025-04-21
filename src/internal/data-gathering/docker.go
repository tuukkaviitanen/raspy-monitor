package datagathering

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"raspy-monitor/src/internal/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/go-units"
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

		if cpuPercentage, err := parsePercentage(DockerStat.CPUPercentage); err != nil {
			log.Printf("Error parsing CPU percentage: %s\n", err)
			continue
		} else {
			fields = append(fields, models.InfluxDbField{
				Name:  "cpu_usage_percentage",
				Value: cpuPercentage,
				Tags: []models.InfluxDbTag{{
					Name: "container_name", Value: DockerStat.Name,
				}},
			})
		}

		if memoryUsedPercentage, err := parsePercentage(DockerStat.MemoryPercentage); err != nil {
			log.Printf("Error parsing Memory percentage: %s\n", err)
			continue
		} else {
			fields = append(fields, models.InfluxDbField{
				Name:  "memory_usage_percentage",
				Value: memoryUsedPercentage,
				Tags: []models.InfluxDbTag{{
					Name: "container_name", Value: DockerStat.Name,
				}},
			})
		}

		if pidCount, err := strconv.Atoi(DockerStat.PIDs); err != nil {
			log.Printf("Error parsing Memory percentage: %s\n", err)
			continue
		} else {
			fields = append(fields, models.InfluxDbField{
				Name:  "pid_count",
				Value: pidCount,
				Tags: []models.InfluxDbTag{{
					Name: "container_name", Value: DockerStat.Name,
				}},
			})
		}
		{
			memoryUsageMatches := doubleStatRegex.FindStringSubmatch(DockerStat.MemoryUsage)

			if memoryUsage, err := units.FromHumanSize(memoryUsageMatches[1]); err != nil {
				log.Printf("Error parsing Memory usage: %s\n", err)
				continue
			} else {
				fields = append(fields, models.InfluxDbField{
					Name:  "memory_usage",
					Value: memoryUsage,
					Tags: []models.InfluxDbTag{{
						Name: "container_name", Value: DockerStat.Name,
					}},
				})
			}
		}
		{
			netIOMatches := doubleStatRegex.FindStringSubmatch(DockerStat.NetIO)

			if dataReceived, err := units.FromHumanSize(netIOMatches[1]); err != nil {
				log.Printf("Error parsing data received: %s\n", err)
				continue
			} else {
				fields = append(fields, models.InfluxDbField{
					Name:  "data_received",
					Value: dataReceived,
					Tags: []models.InfluxDbTag{{
						Name: "container_name", Value: DockerStat.Name,
					}},
				})
			}

			if dataSent, err := units.FromHumanSize(netIOMatches[2]); err != nil {
				log.Printf("Error parsing data sent: %s\n", err)
				continue
			} else {
				fields = append(fields, models.InfluxDbField{
					Name:  "data_sent",
					Value: dataSent,
					Tags: []models.InfluxDbTag{{
						Name: "container_name", Value: DockerStat.Name,
					}},
				})
			}
		}

		{
			blockIOMatches := doubleStatRegex.FindStringSubmatch(DockerStat.BlockIO)

			if dataRead, err := units.FromHumanSize(blockIOMatches[1]); err != nil {
				log.Printf("Error parsing data read: %s\n", err)
				continue
			} else {
				fields = append(fields, models.InfluxDbField{
					Name:  "data_read",
					Value: dataRead,
					Tags: []models.InfluxDbTag{{
						Name: "container_name", Value: DockerStat.Name,
					}},
				})
			}

			if dataWritten, err := units.FromHumanSize(blockIOMatches[2]); err != nil {
				log.Printf("Error parsing data written: %s\n", err)
				continue
			} else {
				fields = append(fields, models.InfluxDbField{
					Name:  "data_written",
					Value: dataWritten,
					Tags: []models.InfluxDbTag{{
						Name: "container_name", Value: DockerStat.Name,
					}},
				})
			}
		}
	}

	return fields
}

var doubleStatRegex = regexp.MustCompile(`^(.*) / (.*)$`)
