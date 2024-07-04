package repositories

import (
	"github.com/gennadyterekhov/gophermart/internal/storage"
)

type Repository struct {
	DB *storage.DB
}

func NewRepository(db *storage.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) Clear() {
	panic("this method is only to satisfy interface. it must not be called")
}
