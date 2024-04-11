package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/models"
	"zatrasz75/Ads_service/pkg/logger"
	"zatrasz75/Ads_service/pkg/mongo"
)

type Store struct {
	l logger.LoggersInterface
	*mongo.Mongo
	cfg *configs.Config
}

func New(mg *mongo.Mongo, l logger.LoggersInterface, cfg *configs.Config) *Store {
	return &Store{l, mg, cfg}
}

// GetSpecificPost Получения конкретного объявления
func (s *Store) GetSpecificPost(id string) (models.Ads, error) {
	// Преобразование строкового ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.l.Error("Не удалось преобразовать строковый ID в ObjectID", err)
		return models.Ads{}, fmt.Errorf("не удалось преобразовать строковый ID в ObjectID: %w", err)
	}

	// Создание фильтра для поиска документа по ID
	filter := bson.M{"_id": objectID}

	// Выполнение поиска документа в коллекции
	var result models.Ads
	err = s.M.Database(s.cfg.Mongo.DbName).Collection(s.cfg.Mongo.CollectionName).FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		s.l.Error("Ошибка при поиске объявления по ID", err)
		return models.Ads{}, fmt.Errorf("ошибка при поиске объявления по ID: %w", err)
	}

	return result, nil
}

// AddPost Добавляет новую запись
func (s *Store) AddPost(ads models.Ads) (string, error) {
	// Создание нового документа для MongoDB
	newAd := bson.M{
		"name":        ads.Name,
		"description": ads.Description,
		"price":       ads.Price,
		"creation":    ads.Creation,
	}

	// Добавление нового документа в коллекцию
	insertResult, err := s.M.Database(s.cfg.Mongo.DbName).Collection(s.cfg.Mongo.CollectionName).InsertOne(context.Background(), newAd)
	if err != nil {
		s.l.Error("Ошибка при добавлении нового объявления: %v", err)
		return "", err
	}

	// Получение ID нового документа
	id := insertResult.InsertedID

	objectID, ok := id.(primitive.ObjectID)
	if !ok {
		s.l.Debug("Не удалось преобразовать ID в ObjectID")
		return "", fmt.Errorf("не удалось преобразовать ID в ObjectID")
	}

	// Преобразование ObjectID в строку
	return objectID.Hex(), nil
}
