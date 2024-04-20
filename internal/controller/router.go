package controller

import (
	"github.com/gin-gonic/gin"
	"zatrasz75/Ads_service/configs"
	_ "zatrasz75/Ads_service/docs"
	"zatrasz75/Ads_service/internal/repository"
	"zatrasz75/Ads_service/pkg/logger"
)

// @title Swagger API
// @version 1.0
// @description ТЗ test_task_backend.
// @description https://github.com/incidentware/test_task_backend/tree/main

// @contact.name Михаил Токмачев
// @contact.url https://t.me/Zatrasz
// @contact.email zatrasz@ya.ru

// @BasePath /

// NewRouter -.
func NewRouter(cfg *configs.Config, l logger.LoggersInterface, repo *repository.Store) *gin.Engine {
	r := gin.Default()
	newEndpoint(r, cfg, l, repo)

	return r
}
