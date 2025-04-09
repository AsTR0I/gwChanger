package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type Host struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type Config struct {
	Hosts          []Host `json:"hosts"`
	TargetHostname string `json:"target_hostname"`
}

func main() {
	configFile := "config.json"
	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Обрабатываем каждый хост
	for i, host := range config.Hosts {
		isSetIp := len(strings.Fields(host.Hostname)) > 0 && host.IP == ""
		if isSetIp {
			resolvedIP := resolveIp(host.Hostname)
			if resolvedIP != "" {
				// Обновляем IP в оригинальной структуре
				config.Hosts[i].IP = resolvedIP
				fmt.Println("Resolved IP for", host.Hostname, ":", resolvedIP)
			} else {
				fmt.Println("Could not resolve IP for:", host.Hostname)
			}
		}
	}

	// Сохраняем изменения в конфигурации
	saveConfig(configFile, config)
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




