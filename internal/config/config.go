package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
}

func LoadConfig() (*Config, error) {
	// Check if the .env file exists.
	if _, err := os.Stat(".env"); err == nil {
		// .env exists; load it.
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
		// Bind environment variables to our config keys.
		viper.AutomaticEnv()
		viper.BindEnv("database.host", "DATABASE_HOST")
		viper.BindEnv("database.port", "DATABASE_PORT")
		viper.BindEnv("database.username", "DATABASE_USERNAME")
		viper.BindEnv("database.password", "DATABASE_PASSWORD")
		viper.BindEnv("database.dbname", "DATABASE_DBNAME")

		// Since environment variables are strings, we might need to convert port.
		if portStr := os.Getenv("DATABASE_PORT"); portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				viper.Set("database.port", port)
			}
		}
	} else {
		// If no .env file exists, load configuration from config.json.
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Error reading config.json: %v", err)
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Error unmarshaling config: %v", err)
		return nil, err
	}
	return &cfg, nil
}
