package mongo

import "time"

// Option -.
type Option func(mongo *Mongo)

func OptionSet(attempts int, timeout time.Duration, dbName, collectionName string) Option {
	return func(m *Mongo) {
		connAttempts(attempts)(m)
		connTimeout(timeout)(m)
		dbNameOption(dbName)(m)
		collectionNameOption(collectionName)(m)
	}
}

// Попытки соединения
func connAttempts(attempts int) Option {
	return func(m *Mongo) {
		m.connAttempts = attempts
	}
}

// Время ожидания соединения
func connTimeout(timeout time.Duration) Option {
	return func(m *Mongo) {
		m.connTimeout = timeout
	}
}

// Имя базы данных
func dbNameOption(dbName string) Option {
	return func(m *Mongo) {
		m.dbName = dbName
	}
}

// Имя коллекции
func collectionNameOption(collectionName string) Option {
	return func(m *Mongo) {
		m.collectionName = collectionName
	}
}
