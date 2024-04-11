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
