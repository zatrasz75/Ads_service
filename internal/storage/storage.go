package storage

import "zatrasz75/Ads_service/models"

type RepositoryInterface interface {
	// AddPost Добавляет новую запись
	AddPost(ads models.Ads) (string, error)
}
