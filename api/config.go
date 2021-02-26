package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Google   ConfigGoogle   `json:"google"`
	Server   ConfigServer   `json:"server"`
	Database ConfigDatabase `json:"database"`
}

type ConfigGoogle struct {
	CredentialsPath   string `json:"credentials_path"`
	DriveRootFolderID string `json:"drive_root_folder_id"`
}

type ConfigServer struct {
	LocalMode bool   `json:"local_mode"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
}

type ConfigDatabase struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
}

func LoadConfig(filePath string) (*Config, error) {
	cfg := &Config{}

	dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	json.Unmarshal(dataBytes, cfg)

	return cfg, nil
}
