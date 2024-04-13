package controller

import (
	"github.com/gorilla/mux"
	"zatrasz75/Ads_service/configs"
	_ "zatrasz75/Ads_service/docs"
	"zatrasz75/Ads_service/internal/repository"
	"zatrasz75/Ads_service/pkg/logger"
)

// @title Swagger API
// @version 1.0
// @description ТЗ test_task_backend.

// @contact.name Михаил Токмачев
// @contact.url https://t.me/Zatrasz
// @contact.email zatrasz@ya.ru

// @BasePath /

// NewRouter -.
func NewRouter(cfg *configs.Config, l logger.LoggersInterface, repo *repository.Store) *mux.Router {
	r := mux.NewRouter()
	newEndpoint(r, cfg, l, repo)

	return r
}
