package config

import (
	"fmt"
	"os"
	"time"

	"go.yaml.in/yaml/v4"
)
type Config struct {
	Database               		DatabaseConfig         		`yaml:"database"`
	Kafka                  		KafkaConfig            		`yaml:"kafka"`
	NotificationServiceSettings NotificationServiceSettings `yaml:"NotificationServiceSettings"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type KafkaConfig struct {
	Host                    string `yaml:"host"`
	Port 					int    `yaml:"port"`
	EventNotificationTopic  string `yaml:"event_notification_topic"`
	NotificationStatusTopic string `yaml:"notification_status_topic"`
}

type NotificationServiceSettings struct {
	NotificationBatchSize int 			`yaml:"notification_batch_size"`
	NotificationTimeout   time.Duration `yaml:"notification_timeout"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}