package repository

import (
	"Tasks/model"
	"context"
	"log"
)

func (r *Repository) GetPeople(ctx context.Context) ([]model.Person, error) {
	var people []model.Person
	err := r.DB.WithContext(ctx).Table("person").Raw("Select id, content, created_at, updated_at, cache From person Where deleted_at is null").Find(&people).Error
	if err != nil {
		log.Printf("GetPeople: %v", err)
		return nil, err
	}
	return people, nil
}
