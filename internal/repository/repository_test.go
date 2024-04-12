package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/models"
	"zatrasz75/Ads_service/pkg/logger"
	"zatrasz75/Ads_service/pkg/mongo"
)

func TestStore_AddPost_GetSpecificPost(t *testing.T) {
	l := logger.NewLogger()

	// Получаем текущий рабочий каталог
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущего рабочего каталога:", err)
		return
	}
	// Построение абсолютного пути к файлу configs.yml
	configPath := filepath.Join(cwd, "..", "..", "configs", "configs.yml")

	// Configuration
	cfg, err := configs.NewConfig(l, configPath)
	if err != nil {
		l.Fatal("ошибка при разборе конфигурационного файла", err)
	}
	mg, err := mongo.New(cfg.Mongo.ConnStr, l, mongo.OptionSet(cfg.Mongo.ConnAttempts, cfg.Mongo.ConnTimeout, cfg.Mongo.DbName, cfg.Mongo.CollectionName))
	if err != nil {
		l.Fatal("нет соединения с базой данных", err)
	}

	repo := New(mg, l, cfg)

	var ct = time.Now().UTC()
	// Создание тестового объявления
	ad := models.Ads{
		Name:        "реклама",
		Description: "Это тестовая реклама",
		Price:       100.05,
		Creation:    ct,
	}

	// Вызов метода AddPost
	id, err := repo.AddPost(ad)
	if err != nil {
		t.Fatalf("Ошибка при добавлении объявления: %v", err)
	}

	// Проверка, что ID возвращается корректно
	if id == "" {
		t.Error("ID объявления не должно быть пустым")
	} else {
		t.Log("OK ID:", id)
	}

	// Получение конкретного объявления по ID
	result, err := repo.GetSpecificPost(id)
	if err != nil {
		t.Fatalf("Ошибка при получении объявления по ID: %v", err)
	}

	// Проверка, что возвращаемые данные соответствуют ожидаемым, с учетом UTC
	if result.Name != ad.Name || result.Description != ad.Description || result.Price != ad.Price {
		t.Errorf("Полученные данные не соответствуют ожидаемым. Ожидалось: %v, Получено: %v", ad, result)
	} else {
		t.Log("OK:", result)
	}
}

func TestStore_GetListPost_AddPost(t *testing.T) {
	l := logger.NewLogger()

	// Получаем текущий рабочий каталог
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущего рабочего каталога:", err)
		return
	}
	// Построение абсолютного пути к файлу configs.yml
	configPath := filepath.Join(cwd, "..", "..", "configs", "configs.yml")

	// Configuration
	cfg, err := configs.NewConfig(l, configPath)
	if err != nil {
		l.Fatal("ошибка при разборе конфигурационного файла", err)
	}
	mg, err := mongo.New(cfg.Mongo.ConnStr, l, mongo.OptionSet(cfg.Mongo.ConnAttempts, cfg.Mongo.ConnTimeout, cfg.Mongo.DbName, cfg.Mongo.CollectionName))
	if err != nil {
		l.Fatal("нет соединения с базой данных", err)
	}

	repo := New(mg, l, cfg)

	var ct = time.Now().UTC()
	ads := []models.Ads{
		{Name: "реклама 1", Description: "Это тестовая реклама 1", Price: 100.05, Creation: ct},
		{Name: "реклама 2", Description: "Это тестовая реклама 2", Price: 200.05, Creation: ct.Add(time.Minute)},
		{Name: "реклама 3", Description: "Это тестовая реклама 3", Price: 300.05, Creation: ct.Add(2 * time.Minute)},
	}

	for _, ad := range ads {
		_, err = repo.AddPost(ad)
		if err != nil {
			t.Fatalf("Ошибка при добавлении объявления: %v", err)
		}
	}

	posts, err := repo.GetListPost(1, "creation", "asc")
	if err != nil {
		t.Fatalf("Ошибка при получении списка объявлений: %v", err)
	}

	// Проверка, что возвращаемые данные соответствуют ожидаемым
	if len(posts) != 0 {
		t.Log("OK:", posts)
	} else {
		t.Errorf("Полученные данные не соответствуют ожидаемым. Ожидалось минимум: %v, Получено: %v", ads, posts)
	}
}
