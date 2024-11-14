package repository

import (
	"Tasks/model"
	"context"
	"log"
)

func (r *Repository) AddPeople(persons []model.Person, ctx context.Context) ([]model.Person, error) {
	var people model.PeopleNoContent
	err := r.DB.WithContext(ctx).Table("person").Create(&persons).Scan(&people).Error
	if err != nil {
		log.Printf("Repository - AddPeople: %v", err)
		return nil, err
	}
	return persons, nil
}
