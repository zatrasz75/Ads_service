package storage

import "zatrasz75/Ads_service/models"

type RepositoryInterface interface {
	// GetListPost Получения списка объявлений
	GetListPost(page int, sortField, sortOrder string) ([]models.Ads, error)
	// GetSpecificPost Получения конкретного объявления
	GetSpecificPost(id string) (models.Ads, error)
	// AddPost Добавляет новую запись
	AddPost(ads models.Ads) (string, error)
}
