package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
	"zatrasz75/Ads_service/pkg/logger"
)

type Config struct {
	Server struct {
		AddrPort     string        `yaml:"port" env:"APP_PORT" env-description:"Server port" env-default:"3131"`
		AddrHost     string        `yaml:"host" env:"APP_IP" env-description:"Server host" env-default:"0.0.0.0"`
		ReadTimeout  time.Duration `yaml:"read-timeout" env:"READ_TIMEOUT" env-description:"Server ReadTimeout" env-default:"3s"`
		WriteTimeout time.Duration `yaml:"write-timeout" env:"WRITE_TIMEOUT" env-description:"Server WriteTimeout" env-default:"3s"`
		IdleTimeout  time.Duration `yaml:"idle-timeout" env:"IDLE_TIMEOUT" env-description:"Server IdleTimeout" env-default:"6s"`
		ShutdownTime time.Duration `yaml:"shutdown-timeout" env:"SHUTDOWN_TIMEOUT" env-description:"Server ShutdownTime" env-default:"10s"`
	} `yaml:"server"`
	Mongo struct {
		ConnStr string `yaml:"connStr" env:"MONGO_CONN_STR" env-description:"MongoDB connection string"`

		Host           string `yaml:"host" env:"MONGO_HOST_DB" env-description:"db host" env-default:"localhost"`
		User           string `yaml:"username" env:"MONGO_USER" env-description:"db username" env-default:"admin"`
		Password       string `yaml:"password" env:"MONGO_PASSWORD" env-description:"db password" env-default:"password"`
		DbName         string `yaml:"db-name" env:"MONGO_DB_NAME" env-description:"db name" env-default:"mongodb"`
		CollectionName string `yaml:"collectionName" env:"MONGO_COLLECTION_NAME" env-description:"collection name" env-default:"ads"`
		Port           string `yaml:"port" env:"MONGO_PORT_DB" env-description:"db port" env-default:"27017"`

		ConnAttempts int           `yaml:"conn-attempts" env:"MONGO_CONN_ATTEMPTS" env-description:"db ConnAttempts" env-default:"5"`
		ConnTimeout  time.Duration `yaml:"conn-timeout" env:"MONGO_TIMEOUT" env-description:"db ConnTimeout" env-default:"1s"`
	} `yaml:"mongo"`
}

func NewConfig(l logger.LoggersInterface, path string) (*Config, error) {
	var cfg Config

	//if err := godotenv.Load(); err != nil {
	//	l.Warn("системе не удается найти указанный файл .env: - %v", err)
	//}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		l.Error("ошибка .env ", err)
		return nil, err
	}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	cfg.Mongo.ConnStr = initDB(cfg)

	return &cfg, nil
}

func initDB(cfg Config) string {
	if cfg.Mongo.ConnStr != "" {
		return cfg.Mongo.ConnStr
	}
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s",
		cfg.Mongo.DbName,
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
	)
}
