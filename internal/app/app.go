package app

import (
	"os"
	"os/signal"
	"syscall"
	"zatrasz75/Ads_service/configs"
	"zatrasz75/Ads_service/internal/controller"
	"zatrasz75/Ads_service/internal/repository"
	"zatrasz75/Ads_service/pkg/logger"
	"zatrasz75/Ads_service/pkg/mongo"
	"zatrasz75/Ads_service/pkg/server"
)

func Run(cfg *configs.Config, l logger.LoggersInterface) {
	mg, err := mongo.New(cfg.Mongo.ConnStr, l, mongo.OptionSet(cfg.Mongo.ConnAttempts, cfg.Mongo.ConnTimeout, cfg.Mongo.DbName, cfg.Mongo.CollectionName))
	if err != nil {
		l.Fatal("нет соединения с базой данных", err)
	}

	repo := repository.New(mg, l, cfg)

	router := controller.NewRouter(cfg, l, repo)

	srv := server.New(router, server.OptionSet(cfg.Server.AddrHost, cfg.Server.AddrPort, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout, cfg.Server.IdleTimeout, cfg.Server.ShutdownTime))

	go func() {
		err = srv.Start()
		if err != nil {
			l.Error("Остановка сервера:", err)
		}
	}()

	l.Info("Запуск сервера на http://" + cfg.Server.AddrHost + ":" + cfg.Server.AddrPort)
	l.Info("Документация Swagger API: http://" + cfg.Server.AddrHost + ":" + cfg.Server.AddrPort + "/swagger/index.html")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("принят сигнал прерывания прерывание %s", s.String())
	case err = <-srv.Notify():
		l.Error("получена ошибка сигнала прерывания сервера", err)
	}

	err = srv.Shutdown()
	if err != nil {
		l.Error("не удалось завершить работу сервера", err)
	}
}
