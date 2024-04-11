package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"zatrasz75/Ads_service/pkg/logger"
)

// Mongo Хранилище данных.
type Mongo struct {
	connAttempts   int
	connTimeout    time.Duration
	dbName         string
	collectionName string

	db *mongo.Client
}

// New Конструктор, принимает строку подключения к БД.
func New(constr string, l logger.LoggersInterface, opts ...Option) (*Mongo, error) {
	m := &Mongo{}

	// Пользовательские параметры
	for _, opt := range opts {
		opt(m)
	}

	mongoOpts := options.Client().ApplyURI(constr)

	var client *mongo.Client
	var err error

	for m.connAttempts > 0 {
		client, err = mongo.Connect(context.Background(), mongoOpts)
		if err == nil {
			// Проверяем, что подключение действительно было установлено
			err = client.Ping(context.Background(), nil)
			if err == nil {
				// Подключение успешно, выходим из цикла
				break
			}
		}
		l.Info("MongoDB пытается подключиться, попыток осталось: %d", m.connAttempts)

		time.Sleep(m.connTimeout)

		m.connAttempts--
	}
	if err != nil {
		return nil, err
	}

	s := Mongo{
		db: client,
	}
	return &s, nil
}
