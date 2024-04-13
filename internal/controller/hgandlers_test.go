package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/internal/repository"
	"zatrasz75/Ads_service/pkg/logger"
	"zatrasz75/Ads_service/pkg/mongo"
)

func Test_api_addPost_getSpecificPost(t *testing.T) {
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

	repo := repository.New(mg, l, cfg)

	a := &api{
		Cfg:  cfg,
		l:    l,
		repo: repo,
	}

	bdy := strings.NewReader(`{
    "name": "заголовок имени",
    "description": "описание",
    "price": 53
    }`)

	req, err := http.NewRequest("POST", "/post", bdy)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.addPost(rr, req)

	// Проверка кода статуса - это то, что мы ожидаем
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Получили code: %v Ожидали %v", status, http.StatusOK)
	} else {
		t.Log("OK:", http.StatusOK, rr.Body.String())
	}

	// Извлечение ID из ответа
	var response struct {
		ID string `json:"id"`
	}
	if err = json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при разборе JSON: %v", err)
	}
	// Проверка полученного ID
	if response.ID == "" {
		t.Error("ID не найден в ответе")
	} else {
		t.Log("Полученный ID:", response.ID)
	}

	reqs, err := http.NewRequest("GET", "/post?id="+response.ID+"&fields=description", nil)
	if err != nil {
		t.Fatal(err)
	}
	rrs := httptest.NewRecorder()

	a.getSpecificPost(rrs, reqs)

	// Проверка кода статуса - это то, что мы ожидаем
	if status := rrs.Code; status != http.StatusOK {
		t.Errorf("Получили code: %v Ожидали %v", status, http.StatusOK)
	} else {
		t.Log("OK:", http.StatusOK, rrs.Body.String())
	}
}

func Test_api_getListPost(t *testing.T) {
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

	repo := repository.New(mg, l, cfg)

	a := &api{
		Cfg:  cfg,
		l:    l,
		repo: repo,
	}

	req, err := http.NewRequest("GET", "/post/list?page=1&sortField=price&sortOrder=asc", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.getListPost(rr, req)

	// Проверка кода статуса - это то, что мы ожидаем
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Получили code: %v Ожидали %v", status, http.StatusOK)
	} else {
		t.Log("OK:", http.StatusOK, rr.Body.String())
	}
}
