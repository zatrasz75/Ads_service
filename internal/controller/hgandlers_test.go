package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

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
	cfg, err := configs.NewConfig(configPath)
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
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Writer = &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}

	a.addPost(c)

	// Проверка кода статуса - это то, что мы ожидаем
	if status := c.Writer.Status(); status != http.StatusOK {
		t.Errorf("Получили code: %v Ожидали %v", status, http.StatusOK)
	} else {
		t.Log("OK:", http.StatusOK, c.Writer.(*responseBodyWriter).body.String())
	}

	// Извлечение ID из ответа
	var response struct {
		ID string `json:"id"`
	}
	if err = json.Unmarshal(c.Writer.(*responseBodyWriter).body.Bytes(), &response); err != nil {
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
	wGet := httptest.NewRecorder()
	cGet, _ := gin.CreateTestContext(wGet)
	cGet.Request = reqs
	cGet.Writer = &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: cGet.Writer}

	a.getSpecificPost(cGet)

	// Проверка кода статуса - это то, что мы ожидаем
	if status := cGet.Writer.Status(); status != http.StatusOK {
		t.Errorf("Получили code: %v Ожидали %v", status, http.StatusOK)
	} else {
		t.Log("OK:", http.StatusOK, cGet.Writer.(*responseBodyWriter).body.String())
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
	cfg, err := configs.NewConfig(configPath)
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
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Writer = &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}

	a.getListPost(c)

	// Проверка кода статуса - это то, что мы ожидаем
	if status := c.Writer.Status(); status != http.StatusOK {
		t.Errorf("Получили code: %v Ожидали %v", status, http.StatusOK)
	} else {
		t.Log("OK:", http.StatusOK, c.Writer.(*responseBodyWriter).body.String())
	}
}
