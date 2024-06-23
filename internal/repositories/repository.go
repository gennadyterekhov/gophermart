package repositories

import "github.com/gennadyterekhov/gophermart/internal/storage"

type Repository struct {
	DB *storage.DB
}

func NewRepository(db *storage.DB) Repository {
	return Repository{
		DB: db,
	}
}
