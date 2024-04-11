package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"time"
	"zatrasz75/Ads_service/pkg/logger"
)

type Config struct {
	Server struct {
		AddrPort     string        `yaml:"port" env:"APP_PORT" env-description:"Server port" env-default:"4141"`
		AddrHost     string        `yaml:"host" env:"APP_IP" env-description:"Server host" env-default:"0.0.0.0"`
		ReadTimeout  time.Duration `yaml:"read-timeout" env:"READ_TIMEOUT" env-description:"Server ReadTimeout" env-default:"3s"`
		WriteTimeout time.Duration `yaml:"write-timeout" env:"WRITE_TIMEOUT" env-description:"Server WriteTimeout" env-default:"3s"`
		IdleTimeout  time.Duration `yaml:"idle-timeout" env:"IDLE_TIMEOUT" env-description:"Server IdleTimeout" env-default:"6s"`
		ShutdownTime time.Duration `yaml:"shutdown-timeout" env:"SHUTDOWN_TIMEOUT" env-description:"Server ShutdownTime" env-default:"10s"`
	} `yaml:"server"`
}

func NewConfig(l logger.LoggersInterface) (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		l.Warn("системе не удается найти указанный файл .env: - %v", err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		l.Error("ошибка .env ", err)
		return nil, err
	}
	if err := cleanenv.ReadConfig("./configs/configs.yml", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
