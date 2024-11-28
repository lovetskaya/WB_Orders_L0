package config

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	ServerAddress   string `json:"server_address"`
	DBConnectionString string `json:"db_connection_string"`
	KafkaBrokers    []string `json:"kafka_brokers"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	file, err := os.Open("service/config/config.json")
	if err != nil {
		return cfg, fmt.Errorf("ошибка при открытии файла: %w", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return cfg, err
}

func NewLogger() (*zap.Logger, error) {
	// Настройка конфигурации логгера
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		OutputPaths:      []string{"logger/app.log"}, // Указываем файл для записи логов
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}

	// Создание логгера
	return config.Build()
}