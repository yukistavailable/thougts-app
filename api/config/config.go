package config

import (
	"log"
	"os"
	"gopkg.in/ini.v1"
)

type ConfigList struct {
	LogFile string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}
	Config = ConfigList{
		LogFile: cfg.Section("thoughts-app").Key("log_file").String(),
	}
}