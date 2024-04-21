package firebaseApp

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/iterator"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/models"
	"zatrasz75/Ads_service/pkg/logger"
)

type FirebaseApp struct {
	l      logger.LoggersInterface
	cfg    *configs.Config
	ctx    context.Context
	app    *firebase.App
	client *firestore.Client
}

func New(l logger.LoggersInterface, cfg *configs.Config, ctx context.Context, app *firebase.App) (*FirebaseApp, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать клиент Firestore: %v", err)
	}

	return &FirebaseApp{l, cfg, ctx, app, client}, nil
}

// GetListPost Получения списка объявлений
func (fa *FirebaseApp) GetListPost(page int, sortField, sortOrder string) ([]models.Ads, error) {
	const pageSize = 10

	// Определение направления сортировки
	var direction firestore.Direction
	if sortOrder == "asc" {
		direction = firestore.Asc
	} else if sortOrder == "desc" {
		direction = firestore.Desc
	} else {
		return nil, fmt.Errorf("некорректное значение sortOrder: %s", sortOrder)
	}

	// Создание запроса с сортировкой
	query := fa.client.Collection("ads").OrderBy(sortField, direction).Limit(pageSize).Offset(int(pageSize * (page - 1)))

	// Выполнение запроса
	iter := query.Documents(fa.ctx)
	defer iter.Stop()

	var posts []models.Ads
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fa.l.Error("Ошибка при получении объявлений", err)
			return nil, fmt.Errorf("ошибка при получении объявлений: %w", err)
		}

		var post models.Ads
		if err := doc.DataTo(&post); err != nil {
			fa.l.Error("Ошибка при декодировании объявления", err)
			return nil, fmt.Errorf("ошибка при декодировании объявления: %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// GetSpecificPost Получения конкретного объявления
func (fa *FirebaseApp) GetSpecificPost(id string) (models.Ads, error) {
	// Создание ссылки на документ в Firestore
	docRef := fa.client.Collection("ads").Doc(id)

	// Получение документа
	doc, err := docRef.Get(fa.ctx)
	if err != nil {
		fa.l.Error("Ошибка при получении объявления по ID", err)
		return models.Ads{}, fmt.Errorf("ошибка при получении объявления по ID: %w", err)
	}

	// Проверка наличия документа
	if !doc.Exists() {
		fa.l.Error("Объявление не найдено", nil)
		return models.Ads{}, fmt.Errorf("объявление не найдено")
	}

	// Декодирование документа в структуру models.Ads
	var post models.Ads
	if err := doc.DataTo(&post); err != nil {
		fa.l.Error("Ошибка при декодировании объявления", err)
		return models.Ads{}, fmt.Errorf("ошибка при декодировании объявления: %w", err)
	}

	return post, nil
}

// AddPost Добавляет новую запись
func (fa *FirebaseApp) AddPost(ads models.Ads) (string, error) {
	// Создание нового документа
	newAd := map[string]interface{}{
		"name":        ads.Name,
		"description": ads.Description,
		"price":       ads.Price,
		"creation":    ads.Creation,
	}

	// Добавление нового документа в коллекцию
	docRef, _, err := fa.client.Collection("ads").Add(fa.ctx, newAd)
	if err != nil {
		fa.l.Error("Ошибка при добавлении нового объявления:", err)
		return "", err
	}

	// Получение ID нового документа
	return docRef.ID, nil
}
