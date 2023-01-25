package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Http struct {
		Port  string `json:"port"`
		Read  int    `json:"read"`
		Write int    `json:"write"`
	} `json:"http"`
	AuthService struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"auth_service"`
}

func InitConfig(path string) (*Config, error) {
	log.Println("init config")
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("init configs: %w", err)
	}
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("init configs: %w", err)
	}
	return &cfg, nil
}
