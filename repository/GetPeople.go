package repository

import (
	"Tasks/model"
	"context"
)

func (r *Repository) GetPeople(ctx context.Context) ([]model.Person, error) {
	var people []model.Person
	err := r.DB.WithContext(ctx).Find(&people).Error
	if err != nil {
		return nil, err
	}
	return people, nil
}
