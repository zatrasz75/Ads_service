package storage

import "zatrasz75/Ads_service/models"

type RepositoryInterface interface {
	// GetSpecificPost Получения конкретного объявления
	GetSpecificPost(id string) (models.Ads, error)
	// AddPost Добавляет новую запись
	AddPost(ads models.Ads) (string, error)
}
