package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"path/filepath"
	"time"
	"net"
)

const (
	Version        = "24.10.9"
	RN             = "\r\n"
	R              = "\r"
	N              = "\n"
)

var (
	hostsFile     = "/etc/hosts"
	bestTime      = 999999.0
	bestHost      = ""
	bestIP        = ""
	logsDir       string
	LogSaveDays   = 10
)

// colors
var (
	colorWhite            = "\033[97m"
	colorReset            = "\033[0m"
	colorCyanBackground   = "\033[46m"
	colorWhiteBackground  = "\033[107m"
)

type Host struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type Config struct {
	Hosts          []Host `json:"hosts"`
	TargetHostname string  `json:"target_hostname"`
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Can't get current dir")
		return
	}

	logsDir = filepath.Join(dir, "logs")

	// Create logs directory if it doesn't exist
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err = os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating logs directory:", err)
			return
		}
	}

	logFileName := fmt.Sprintf("%s/gw_changer_log_%s.log", logsDir, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	writeLog(file, "Program started")

	// Load hosts configuration
	config, err := loadConfig("./config.json")
	if err != nil {
		writeLog(file, fmt.Sprintf("Error loading hosts: %v", err))
		return
	}

	// Resolve IPs for hosts with empty IP
	for i, host := range config.Hosts {
		if host.IP == "" {
			resolvedIP := resolveIP(host.Hostname)
			if resolvedIP != "" {
				config.Hosts[i].IP = resolvedIP
				writeLog(file, fmt.Sprintf("Resolved IP for %s: %s", host.Hostname, resolvedIP))
			} else {
				writeLog(file, fmt.Sprintf("Could not resolve IP for %s", host.Hostname))
			}
		}
	}

	// Save updated configuration
	saveConfig("config.json", config)

	// Process command line arguments
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "--help":
			printHelpMessage()
			return
		case "-v", "--version":
			fmt.Printf("gwChanger v%s\n", Version)
			return
		}
	}

	if len(os.Args) > 1 {
		if value := myFlagGetValue("-lsd"); value != "" {
			localLogSaveDays, err := strconv.Atoi(value)
			if err == nil {
				LogSaveDays = localLogSaveDays
			} else {
				LogSaveDays = 10
			}
		}
	}

	deleteOldLogs(logsDir, LogSaveDays)

	// Find the best host based on elapsed time
	var bestHostFound bool

	for _, host := range config.Hosts {
		hostArgs := strings.Fields(host.Hostname)
		cmdArgs := append([]string{"o"}, hostArgs...)
		cmdArgs = append(cmdArgs, "-pt")

		cmd := exec.Command("sipc", cmdArgs...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			writeLog(file, fmt.Sprintf("Error executing sipc for host %s: %v", host.Hostname, err))
			continue
		}

		outputStr := string(output)
		writeLog(file, fmt.Sprintf("Output from sipc for host %s: %s", host.Hostname, outputStr))

		if strings.Contains(outputStr, "Elapsed time") {
			lines := strings.Split(outputStr, "\n")
			var elapsedTimeStr string
			for _, line := range lines {
				if strings.Contains(line, "Elapsed time") {
					parts := strings.Fields(line)
					if len(parts) >= 3 {
						elapsedTimeStr = parts[2]
						elapsedTime, err := strconv.ParseFloat(elapsedTimeStr, 64)
						if err != nil {
							writeLog(file, fmt.Sprintf("Error parsing elapsed time: %v", err))
							return
						}
						if bestTime > elapsedTime {
							bestTime = elapsedTime
							bestHost = host.Hostname
							bestIP = host.IP
							bestHostFound = true
							logMessage := fmt.Sprintf("Best time: %f for host: %s with IP: %s", bestTime, bestHost, bestIP)
							writeLog(file, logMessage)
						}
					}
				}
			}
		}
	}

	// Update /etc/hosts if the best host was found
	if bestHostFound {
		changeEtcHosts(file, bestIP, config.TargetHostname)
	}
}

// Load configuration from a JSON file
func loadConfig(filename string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, nil
}

// Save configuration to a JSON file
func saveConfig(filePath string, config Config) {
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(config); err != nil {
		// Handle encoding error
	}
}

// Update /etc/hosts with the new IP for the target hostname
func changeEtcHosts(file *os.File, newIP string, targetHostname string) {
	data, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		writeLog(file, fmt.Sprintf("Error reading %s: %v", hostsFile, err))
		return
	}

	lines := strings.Split(string(data), "\n")
	var found bool
	var updatedLines []string

	for _, line := range lines {
		if strings.Contains(line, targetHostname) {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				if parts[0] == newIP {
					writeLog(file, fmt.Sprintf("IP %s already set for %s, no changes made.", newIP, targetHostname))
					return // IP already set, no changes
				} else {
					line = fmt.Sprintf("%s %s", newIP, targetHostname) // Update IP
					found = true
					writeLog(file, fmt.Sprintf("Updated %s to new IP: %s", targetHostname, newIP))
				}
			}
		}
		updatedLines = append(updatedLines, line)
	}

	if !found {
		updatedLines = append(updatedLines, fmt.Sprintf("%s %s", newIP, targetHostname))
		writeLog(file, fmt.Sprintf("Added new entry: %s %s", newIP, targetHostname))
	}

	err = ioutil.WriteFile(hostsFile, []byte(strings.Join(updatedLines, "\n")), 0644)
	if err != nil {
		writeLog(file, fmt.Sprintf("Error writing to %s: %v", hostsFile, err))
	}
}

// Log a message to the log file
func writeLog(file *os.File, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("%s: %s%s", timestamp, message, N)
	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

// Delete old log files
func deleteOldLogs(logsDir string, days int) {
	files, err := ioutil.ReadDir(logsDir)
	if err != nil {
		return
	}

	expiration := time.Now().AddDate(0, 0, -days)

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "gw_changer_log_") || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		dateStr := strings.TrimPrefix(strings.TrimSuffix(file.Name(), ".log"), "gw_changer_log_")
		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if fileDate.Before(expiration) {
			fullPath := filepath.Join(logsDir, file.Name())
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Println("Error deleting old log file:", err)
			}
		}
	}
}

// Get the value for a flag from command line arguments
func myFlagGetValue(flag string) string {
	for i, arg := range os.Args {
		if arg == flag && len(os.Args) > i+1 {
			return os.Args[i+1]
		}
	}
	return ""
}

// Resolve the IP address for a given hostname
func resolveIP(hostname string) string {
	addrs, err := net.LookupHost(hostname)
	if err != nil || len(addrs) == 0 {
		return ""
	}
	return addrs[0]
}

// Print help message
func printHelpMessage() {
	fmt.Println("Usage:")
	fmt.Println("  gwChanger [options]")
	fmt.Println("Options:")
	fmt.Println("  -h, --help           Show help message")
	fmt.Println("  -v, --version        Show version")
	fmt.Println("  -lsd N               Set log save days (default 10)")
}
