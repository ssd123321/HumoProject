package repository

import (
	"Tasks/model"
	"context"
	"log"
	"time"
)

func (r *Repository) AddPerson(person *model.Person, ctx context.Context) (*model.Person, error) {
	err := r.DB.WithContext(ctx).Table("person").Create(person).Scan(&person).Error
	if err != nil {
		log.Printf("Failed to AddPerson: %v", err)
		return nil, err
	}
	return person, nil
}
func (r *Repository) GetPersonByLogin(login string) (*model.Person, error) {
	var person model.Person
	err := r.DB.Raw("Select * From person Where login = ?", login).First(&person).Error
	if err != nil {
		log.Printf("Repository - GetPersonByLogin: %v", err)
		return nil, err
	}
	return &person, nil
}
func (r *Repository) AddPeople(persons []model.Person, ctx context.Context) ([]model.Person, error) {
	var people model.PeopleNoContent
	err := r.DB.WithContext(ctx).Table("person").Create(&persons).Scan(&people).Error
	if err != nil {
		log.Printf("Repository - AddPeople: %v", err)
		return nil, err
	}
	return persons, nil
}
func (r *Repository) DeletePerson(ctx context.Context, time1 *time.Time) (int, error) {
	var id int
	err := r.DB.WithContext(ctx).Exec("UPDATE person set deleted_at = ? Where id = ? and deleted_at is null", time.Now(), ctx.Value("person_id").(int)).Error
	if err != nil {
		log.Printf("Failed to DeletePerson: %v", err)
		return 0, err
	}
	id = ctx.Value("person_id").(int)
	return id, nil
}
func (r *Repository) GetPeople(ctx context.Context) ([]model.Person, error) {
	var people []model.Person
	err := r.DB.WithContext(ctx).Table("person").Raw("Select id, content, created_at, updated_at, cache From person Where deleted_at is null").Find(&people).Error
	if err != nil {
		log.Printf("GetPeople: %v", err)
		return nil, err
	}
	return people, nil
}
func (r *Repository) GetPersonByID(ctx context.Context) (*model.PersonNew, error) {
	var person model.PersonNew
	err := r.DB.Table("person").First(&person, ctx.Value("id").(int)).Error
	if err != nil {
		log.Printf("Failed to GetPerson - r.DB.Raw(1): %v", err)
		return nil, err
	}
	log.Printf("+%v", person)
	return &person, nil
}
func (r *Repository) UpdatePerson(person *model.Person, ctx context.Context) error {
	err := r.DB.WithContext(ctx).Exec("UPDATE Person Set content = ?  Where id = ? and person.deleted_at is null returning content::text", person.Content, ctx.Value("id").(int)).Error
	if err != nil {
		return err
	}
	return nil
}
