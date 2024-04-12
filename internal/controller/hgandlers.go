package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"strconv"
	"time"
	"zatrasz75/Ads_service/configs"
	_ "zatrasz75/Ads_service/docs"
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
	r.HandleFunc("/post/list", en.getListPost).Methods(http.MethodGet)
	r.HandleFunc("/post", en.getSpecificPost).Methods(http.MethodGet)
	r.HandleFunc("/post", en.addPost).Methods(http.MethodPost)

	r.HandleFunc("/", en.home).Methods(http.MethodGet)

	// Swagger UI
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

// @Summary Получение списка объявлений
// @Description Метод для получения списка объявлений с возможностью сортировки по цене или дате создания, а также пагинации.
// @Description Возвращает список объявлений с указанными параметрами сортировки и пагинации.
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы для пагинации (по умолчанию 1)"
// @Param sortField query string false "Поле для сортировки (например, creation или price)"
// @Param sortOrder query string false "Порядок сортировки (asc или desc)"
// @Success 200 {array} models.Response
// @Failure 500 {string} string "Ошибка при получении списка объявлений"
// @Failure 500 {string} string "Ошибка при сериализации списка объявлений в JSON"
// @Router /post/list [get]
// @OperationId getListPost
func (a *api) getListPost(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	pageStr := queryParams.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	sortField := queryParams.Get("sortField")
	sortOrder := queryParams.Get("sortOrder")

	ads, err := a.repo.GetListPost(page, sortField, sortOrder)
	if err != nil {
		http.Error(w, "Ошибка при получении списка объявлений", http.StatusInternalServerError)
		a.l.Error("Ошибка при получении списка объявлений", err)
		return
	}

	// Создание среза для хранения только необходимых полей
	var response []struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	// Заполнение среза данными из ads
	for _, ad := range ads {
		response = append(response, struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		}{
			Name:  ad.Name,
			Price: ad.Price,
		})
	}

	// Установка заголовка Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Ошибка при сериализации списка объявлений в JSON", http.StatusInternalServerError)
		a.l.Error("Ошибка при сериализации списка объявлений в JSON", err)
		return
	}
}

// @Summary Получение конкретного объявления по ID
// @Description Метод для получения информации о конкретном объявлении по его уникальному идентификатору.
// @Description Возвращает данные объявления, включая название, описание и цену.
// @Description Если поля название объявления или цена отсутствуют возвращает ошибку 400
// @Description Если запрошен параметр "fields" со значением "description", возвращает также описание объявления.
// @Accept json
// @Produce json
// @Param id query string true "ID объявления"
// @Param fields query string false "Опциональные поля для запроса (например, description)"
// @Success 200 {object} models.Response
// @Failure 400 {string} string "Не удалось получить параметр id"
// @Failure 400 {string} string "Обязательные поля объявления отсутствуют"
// @Failure 500 {string} string "Ошибка при получении данных"
// @Router /post [get]
// @OperationId getSpecificPost
func (a *api) getSpecificPost(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	idStr := queryParams.Get("id")
	if idStr == "" {
		a.l.Debug("Не удалось получить параметр id")
		http.Error(w, "Не удалось получить параметр id", http.StatusBadRequest)
		return
	}

	ads, err := a.repo.GetSpecificPost(idStr)
	if err != nil {
		a.l.Error("Ошибка при получении данных", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	// Проверка наличия обязательных полей
	if ads.Name == "" || ads.Price == 0 {
		a.l.Debug("Обязательные поля объявления отсутствуют")
		http.Error(w, "Обязательные поля объявления отсутствуют", http.StatusBadRequest)
		return
	}

	// Проверка наличия параметра fields для запроса опциональных полей
	fields := queryParams.Get("fields")
	if fields == "description" {
		response := struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
		}{
			Name:        ads.Name,
			Description: ads.Description,
			Price:       ads.Price,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			a.l.Error("Ошибка при сериализации ответа JSON", err)
			http.Error(w, "Ошибка при сериализации ответа JSON", http.StatusInternalServerError)
			return
		}
	} else {
		// Если не запрошено описание, возвращаем название и цену
		response := struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		}{
			Name:  ads.Name,
			Price: ads.Price,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			a.l.Error("Ошибка при сериализации ответа JSON", err)
			http.Error(w, "Ошибка при сериализации ответа JSON", http.StatusInternalServerError)
			return
		}
	}

}

// @Summary Создание нового объявления
// @Description Метод для добавления нового объявления в систему.
// @Description Принимает поля: название, описание, цена (name , description , price).
// @Description Обязательные поля: название и цена (name и price).
// @Description Возвращает ID созданного объявления и код результата (ошибка или успех).
// @Accept json
// @Produce json
// @Param ads body models.Ads true "Объявление"
// @Success 200 {object} models.Response
// @Failure 400 {string} string "не удалось проанализировать запрос JSON"
// @Failure 400 {string} string "Обязательные поля name или price объявления отсутствуют"
// @Failure 500 {string} string "не удалось округлить цену"
// @Failure 500 {string} string "Ошибка при добавлении данных"
// @Failure 500 {string} string "не удалось сериализовать ответ JSON"
// @Router /post [post]
// @OperationId addPost
func (a *api) addPost(w http.ResponseWriter, r *http.Request) {
	var p models.Ads

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "не удалось проанализировать запрос JSON", http.StatusBadRequest)
		a.l.Error("не удалось проанализировать запрос JSON", err)
		return
	}
	if p.Name == "" || p.Price == 0 {
		a.l.Debug("Обязательные поля name или price объявления отсутствуют")
		http.Error(w, "Обязательные поля name или price объявления отсутствуют", http.StatusBadRequest)
		return
	}
	p.Creation = time.Now()

	// Округление Price до двух знаков после запятой
	roundedPriceStr := fmt.Sprintf("%.2f", p.Price)
	roundedPrice, err := strconv.ParseFloat(roundedPriceStr, 64)
	if err != nil {
		http.Error(w, "не удалось округлить цену", http.StatusInternalServerError)
		a.l.Error("не удалось округлить цену", err)
		return
	}
	p.Price = roundedPrice

	id, err := a.repo.AddPost(p)
	if err != nil {
		a.l.Error("Ошибка при добавлении данных", err)
		http.Error(w, "Ошибка при добавлении данных", http.StatusInternalServerError)
		return
	}
	response := models.Response{
		ID: id,
	}
	//response := struct {
	//	ID string
	//}{
	//	ID: id,
	//}

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
