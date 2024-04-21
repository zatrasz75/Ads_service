package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"log"
	"zatrasz75/Ads_service/configs"
)

func NewFirebase(cfg *configs.Config) (*firebase.App, error) {
	opt := option.WithCredentialsFile(cfg.Firebase.ConnStr)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Не удалось создать приложение Firebase: %v", err)
	}

	return app, nil
}
