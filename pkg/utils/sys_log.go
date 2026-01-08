package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"uni-search-hub/pkg/common"

	"github.com/gin-gonic/gin"
)

func SysLog(s string) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), s)
}

func SysError(s string) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), s)
}

func FatalLog(v ...any) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[FATAL] %v | %v \n", t.Format("2006/01/02 - 15:04:05"), v)
	os.Exit(1)
}

func LogStartupSuccess(startTime time.Time, port string) {

	duration := time.Since(startTime)
	durationMs := duration.Milliseconds()

	// Get network IPs
	networkIps := GetNetworkIps()

	// Print blank line for spacing
	fmt.Fprintf(gin.DefaultWriter, "\n")

	// Print the main success message
	fmt.Fprintf(gin.DefaultWriter, "  \033[32m%s %s\033[0m  ready in %d ms\n", common.SystemName, common.Version, durationMs)
	fmt.Fprintf(gin.DefaultWriter, "\n")

	// Skip fancy startup message in container environments
	if !IsRunningInContainer() {
		// Print local URL
		fmt.Fprintf(gin.DefaultWriter, "  ➜  \033[1mLocal:\033[0m   http://localhost:%s/\n", port)
	}

	// Print network URLs
	for _, ip := range networkIps {
		fmt.Fprintf(gin.DefaultWriter, "  ➜  \033[1mNetwork:\033[0m http://%s:%s/\n", ip, port)
	}

	// Print blank line for spacing
	fmt.Fprintf(gin.DefaultWriter, "\n")
}

// IsRunningInContainer detects if the application is running inside a container
func IsRunningInContainer() bool {
	// Method 1: Check for .dockerenv file (Docker containers)
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Method 2: Check cgroup for container indicators
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := string(data)
		if strings.Contains(content, "docker") ||
			strings.Contains(content, "containerd") ||
			strings.Contains(content, "kubepods") ||
			strings.Contains(content, "/lxc/") {
			return true
		}
	}

	// Method 3: Check environment variables commonly set by container runtimes
	containerEnvVars := []string{
		"KUBERNETES_SERVICE_HOST",
		"DOCKER_CONTAINER",
		"container",
	}

	for _, envVar := range containerEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	// Method 4: Check if init process is not the traditional init
	if data, err := os.ReadFile("/proc/1/comm"); err == nil {
		comm := strings.TrimSpace(string(data))
		// In containers, process 1 is often not "init" or "systemd"
		if comm != "init" && comm != "systemd" {
			// Additional check: if it's a common container entrypoint
			if strings.Contains(comm, "docker") ||
				strings.Contains(comm, "containerd") ||
				strings.Contains(comm, "runc") {
				return true
			}
		}
	}

	return false
}

func GetNetworkIps() []string {
	var networkIps []string
	ips, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return networkIps
	}

	for _, a := range ips {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip := ipNet.IP.String()
				// Include common private network ranges
				if strings.HasPrefix(ip, "10.") ||
					strings.HasPrefix(ip, "172.") ||
					strings.HasPrefix(ip, "192.168.") {
					networkIps = append(networkIps, ip)
				}
			}
		}
	}
	return networkIps
}
