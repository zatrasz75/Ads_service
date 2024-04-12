package main

import (
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/internal/app"
	"zatrasz75/Ads_service/pkg/logger"
)

func main() {
	l := logger.NewLogger()

	// Configuration
	cfg, err := configs.NewConfig(l)
	if err != nil {
		l.Fatal("ошибка при разборе конфигурационного файла", err)
	}

	app.Run(cfg, l)
}
