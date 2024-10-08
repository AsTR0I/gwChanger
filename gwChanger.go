package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Host struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type Config struct {
	Hosts          []Host `json:"hosts"`
	TargetHostname string  `json:"target_hostname"` // Добавляем поле для целевого хоста
}

var (
	hostsFile = "/etc/hosts"
	bestTime  = 999999.0
	bestHost  = ""
	bestIP    = ""
	logsDir   = "logs"
)

func main() {
	// Создаём директорию, если её не существует
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		return
	}

	logFileName := fmt.Sprintf("%s/gw_changer_log_%s.log", logsDir, time.Now().Format("2006-01-02T15-04-05"))
	// Создание файла с логом с датой
	file, err := os.Create(logFileName)
	if err != nil {
		return
	}
	defer file.Close()

	// Начало логов
	writeLog(file, "Program started")

	// Загрузка конфигурации хостов
	config, err := loadConfig("config.json") // Изменяем на loadConfig
	if err != nil {
		writeLog(file, fmt.Sprintf("Error loading hosts: %v", err))
		return
	}

	// Объявляем переменные для отслеживания лучшего хоста
	var bestHostFound bool

	for _, host := range config.Hosts {
		hostArgs := strings.Fields(host.Hostname)

		cmdArgs := append([]string{"o"}, hostArgs...)
		cmdArgs = append(cmdArgs, "-pt")

		cmd := exec.Command("sipc", cmdArgs...)

		// Получаем вывод и код программы
		output, err := cmd.CombinedOutput()

		// Если код не 0, останавливаем
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
	// Вызываем функцию changeEtcHosts только если был найден лучший хост
	if bestHostFound {
		changeEtcHosts(file, bestIP, config.TargetHostname)
	}
}

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

func changeEtcHosts(file *os.File, newIP string, targetHostname string) {
	// Читаем файл /etc/hosts
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
					return // IP уже установлен, ничего не делаем
				} else {
					line = fmt.Sprintf("%s %s", newIP, targetHostname) // Заменяем IP
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

	// Записываем обновлённые строки обратно в файл
	err = ioutil.WriteFile(hostsFile, []byte(strings.Join(updatedLines, "\n")), 0644)
	if err != nil {
		writeLog(file, fmt.Sprintf("Error writing to %s: %v", hostsFile, err))
	}
}

// Функция для записи логов в файл
func writeLog(file *os.File, message string) {
	logEntry := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	file.WriteString(logEntry)
}
