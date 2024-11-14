package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) NewRepository(Db *gorm.DB) *Repository {
	return &Repository{
		DB: Db,
	}
}
