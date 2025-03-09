package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Projects map[string]string `json:"projects"`
}

func LoadConfig() Config {
	var cfg Config

	// Check if .env exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading .env file: %v", err)
		} else {
			projects := make(map[string]string)
			if web1 := os.Getenv("WEB1_URL"); web1 != "" {
				projects["web1"] = web1
			}
			if mobile := os.Getenv("MOBILE_URL"); mobile != "" {
				projects["mobile1"] = mobile
			}
			if len(projects) > 0 {
				cfg.Projects = projects
				return cfg
			}
		}
	}

	// Fallback to config.json
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Printf("Error reading config.json: %v", err)
		// Provide defaults
		cfg.Projects = map[string]string{
			"web1":    "http://localhost:5000",
			"mobile1": "http://localhost:5001",
		}
		return cfg
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("Error parsing config.json: %v", err)
		cfg.Projects = map[string]string{
			"web1":    "http://localhost:5000",
			"mobile1": "http://localhost:5001",
		}
	}
	return cfg
}
