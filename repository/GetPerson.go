package repository

import (
	"Tasks/model"
	"context"
	"encoding/json"
	"log"
)

func (r *Repository) GetPersonByID(ctx context.Context) (*model.Person, error) {
	var person model.DBPerson
	err := r.DB.Raw("Select id, content::text, deleted_at, created_at, updated_at From person where id = ? and deleted_at is null", ctx.Value("id")).First(&person).Error
	if err != nil {
		log.Printf("Failed to GetPerson - r.DB.Raw(1): %v", err)
		return nil, err
	}
	var content model.Content
	err = json.Unmarshal([]byte(person.Content), &content)
	if err != nil {
		return nil, err
	}
	return &model.Person{
		ID: person.ID,
		Content: model.Content{
			Name:      content.Name,
			Age:       content.Age,
			Dimension: content.Dimension,
		},
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
		DeletedAt: person.DeletedAt,
		Cache:     false,
	}, nil
}
