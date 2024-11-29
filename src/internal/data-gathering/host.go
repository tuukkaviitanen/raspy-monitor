package datagathering

import (
	"fmt"
	"raspy-monitor/src/internal/models"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

func GetHostData() (*models.HostInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("Error getting host info: %v\n", err)
	}

	parsedHostInfo := &models.HostInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformFamily:  hostInfo.PlatformFamily,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		KernelArch:      hostInfo.KernelArch,
		BootTime:        time.Unix(int64(hostInfo.BootTime), 0),
	}

	return parsedHostInfo, nil
}
