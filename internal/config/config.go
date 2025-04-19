package config

import (
	"fmt"
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
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
}

// MinIOConfig holds the MinIO specific configuration
type MinIOConfig struct {
	Endpoint   string `mapstructure:"endpoint"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	UseSSL     bool   `mapstructure:"usessl"`
	BucketName string `mapstructure:"bucket"`
}

// RabbitMQConfig holds the RabbitMQ specific configuration.
type RabbitMQConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	ExchangeName string `mapstructure:"exchange"`
}

// AMQPURL builds the AMQP URL from the individual RabbitMQ configuration fields.
func (r *RabbitMQConfig) AMQPURL() string {
	fmt.Println("r", r)
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.Username, r.Password, r.Host, r.Port)
}

func LoadConfig() (*Config, error) {
	// Check if the .env file exists.
	if _, err := os.Stat(".env"); err == nil {
		fmt.Println("loading env")
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
		viper.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
		viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
		viper.BindEnv("rabbitmq.username", "RABBITMQ_USERNAME")
		viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")
		viper.BindEnv("rabbitmq.exchange", "RABBITMQ_EXCHANGE_NAME")
		viper.BindEnv("minio.endpoint", "MINIO_ENDPOINT")
		viper.BindEnv("minio.username", "MINIO_USERNAME")
		viper.BindEnv("minio.password", "MINIO_PASSWORD")
		viper.BindEnv("minio.usessl", "MINIO_USE_SSL")
		viper.BindEnv("minio.bucket", "MINIO_BUCKET")
		//Since environment variables are strings, we might need to convert port.
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
	fmt.Println(cfg)
	return &cfg, nil
}
