package main

import (
	"fmt"
	"os"
	"path/filepath"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/internal/app"
	"zatrasz75/Ads_service/pkg/logger"
)

func main() {
	l := logger.NewLogger()

	// Получаем текущий рабочий каталог
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущего рабочего каталога:", err)
		return
	}
	// Построение абсолютного пути к файлу configs.yml
	configPath := filepath.Join(cwd, "configs", "configs.yml")

	// Configuration
	cfg, err := configs.NewConfig(l, configPath)
	if err != nil {
		l.Fatal("ошибка при разборе конфигурационного файла", err)
	}

	app.Run(cfg, l)
}
