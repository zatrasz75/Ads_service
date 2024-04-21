package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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

func newEndpoint(r *gin.Engine, cfg *configs.Config, l logger.LoggersInterface, repo *repository.Store) {
	en := &api{cfg, l, repo}
	r.GET("/posts/list", en.getListPost)
	r.GET("/posts", en.getSpecificPost)
	r.POST("/posts", en.addPost)

	r.GET("/", en.home)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
// @Router /posts/list [get]
// @OperationId getListPost
func (a *api) getListPost(c *gin.Context) {
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	sortField := c.Query("sortField")
	sortOrder := c.Query("sortOrder")

	ads, err := a.repo.GetListPost(page, sortField, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "при получении списка объявлений"})
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

	c.JSON(http.StatusOK, response)
	if len(c.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации списка объявлений в JSON"})
		a.l.Debug("Ошибка при сериализации списка объявлений в JSON")
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
// @Router /posts [get]
// @OperationId getSpecificPost
func (a *api) getSpecificPost(c *gin.Context) {
	idStr := c.Query("id")

	if idStr == "" {
		a.l.Debug("не удалось получить параметр id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось получить параметр id"})

		return
	}

	ads, err := a.repo.GetSpecificPost(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении данных"})
		a.l.Error("ошибка при получении данных", err)
		return
	}

	// Проверка наличия обязательных полей
	if ads.Name == "" || ads.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "обязательные поля объявления отсутствуют"})
		a.l.Debug("обязательные поля объявления отсутствуют")
		return
	}

	// Проверка наличия параметра fields для запроса опциональных полей
	fields := c.Query("fields")
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

		c.JSON(http.StatusOK, response)

	} else {
		// Если не запрошено описание, возвращаем название и цену
		response := struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		}{
			Name:  ads.Name,
			Price: ads.Price,
		}

		c.JSON(http.StatusOK, response)
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации списка объявлений в JSON"})
			a.l.Debug("Ошибка при сериализации списка объявлений в JSON")
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
// @Router /posts [post]
// @OperationId addPost
func (a *api) addPost(c *gin.Context) {
	var p models.Ads

	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось проанализировать запрос JSON"})
		a.l.Error("не удалось проанализировать запрос JSON", err)
		return
	}
	if p.Name == "" || p.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "обязательные поля name или price объявления отсутствуют"})
		a.l.Debug("обязательные поля name или price объявления отсутствуют")
		return
	}
	p.Creation = time.Now()

	// Округление Price до двух знаков после запятой
	roundedPriceStr := fmt.Sprintf("%.2f", p.Price)
	roundedPrice, err := strconv.ParseFloat(roundedPriceStr, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось округлить цену"})
		a.l.Error("не удалось округлить цену", err)
		return
	}
	p.Price = roundedPrice

	id, err := a.repo.AddPost(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при добавлении данных"})
		a.l.Error("ошибка при добавлении данных", err)
		return
	}
	response := models.Response{
		ID: id,
	}

	c.JSON(http.StatusOK, response)
	if len(c.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации списка объявлений в JSON"})
		a.l.Debug("Ошибка при сериализации списка объявлений в JSON")
	}
}

func (a *api) home(c *gin.Context) {
	// Выводим дополнительную строку на страницу
	str := "Добро пожаловать! "

	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(str))
	if len(c.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации списка объявлений в JSON"})
		a.l.Debug("Ошибка при сериализации списка объявлений в JSON")
	}
}
