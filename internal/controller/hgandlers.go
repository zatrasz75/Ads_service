package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/internal/repository"
	"zatrasz75/Ads_service/internal/storage"
	"zatrasz75/Ads_service/models"
	"zatrasz75/Ads_service/pkg/logger"
)

type api struct {
	Cfg  *configs.Config
	l    logger.LoggersInterface
	repo storage.RepositoryInterface
}

func newEndpoint(r *mux.Router, cfg *configs.Config, l logger.LoggersInterface, repo *repository.Store) {
	en := &api{cfg, l, repo}
	r.HandleFunc("/post", en.addPost).Methods(http.MethodPost)

	r.HandleFunc("/", en.home).Methods(http.MethodGet)

	// Swagger UI
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))
	//r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

// Метод создания объявления
func (a *api) addPost(w http.ResponseWriter, r *http.Request) {
	var p models.Ads

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "не удалось проанализировать запрос JSON", http.StatusBadRequest)
		a.l.Error("не удалось проанализировать запрос JSON", err)
		return
	}
	p.Creation = time.Now()
	fmt.Println(p)
	id, err := a.repo.AddPost(p)
	if err != nil {
		a.l.Error("Ошибка при добавлении данных", err)
		http.Error(w, "Ошибка при добавлении данных", http.StatusInternalServerError)
		return
	}
	fmt.Println("id:", id)

	response := struct {
		ID string
	}{
		ID: id,
	}

	// Установка заголовка Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Сериализация структуры ответа в JSON и запись в http.ResponseWriter
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "не удалось сериализовать ответ JSON", http.StatusInternalServerError)
		a.l.Error("не удалось сериализовать ответ JSON", err)
		return
	}
}

func (a *api) home(w http.ResponseWriter, _ *http.Request) {
	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Выводим дополнительную строку на страницу
	str := []byte("Добро пожаловать! ")

	_, err := fmt.Fprintf(w, "<p>%s</p>", str)
	if err != nil {
		http.Error(w, "Ошибка записи на страницу", http.StatusInternalServerError)
		a.l.Error("Ошибка записи на страницу", err)
	}
}
