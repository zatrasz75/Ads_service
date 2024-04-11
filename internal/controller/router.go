package controller

import (
	"github.com/gorilla/mux"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/pkg/logger"
)

// NewRouter -.
func NewRouter(cfg *configs.Config, l logger.LoggersInterface) *mux.Router {
	r := mux.NewRouter()
	//newEndpoint(r, cfg, l, repo)
	return r
}
