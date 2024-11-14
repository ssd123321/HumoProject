package repository

import (
	"Tasks/model"
	"context"
	"log"
)

func (r *Repository) AddPerson(person *model.Person, ctx context.Context) (*model.Person, error) {
	err := r.DB.WithContext(ctx).Table("person").Create(&person).Error
	if err != nil {
		log.Printf("Failed to AddPerson: %v", err)
		return nil, err
	}
	return person, nil
}
