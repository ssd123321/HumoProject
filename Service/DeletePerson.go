package Service

import (
	"Tasks/model"
	"Tasks/time"
	"Tasks/utils"
	"context"
	"errors"
	"fmt"
	"log"
)

func (s *Service) Login(person *model.SigningRequest) (string, string, error) {
	p, err := s.repository.GetPersonByLogin(person.Login)
	if err != nil {
		return "", "", err
	}
	b := utils.CheckPasswordHash(person.Password, p.Password)
	if !b {
		return "", "", fmt.Errorf("password is incorrect")
	}
	refreshJWT, err := utils.GenerateRefreshJWT(p)
	if err != nil {
		return "", "", err
	}
	accessJWT, err := utils.GenerateAccessJWT(refreshJWT)
	err = s.repository.AddRefresh(p.ID, refreshJWT)
	if err != nil {
		return "", "", err
	}
	return refreshJWT, accessJWT, nil
}

func (s *Service) DeletePerson(ctx context.Context) (int, error) {
	d, err := time.GetCurrentDate()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	person, err := s.repository.GetPersonByID(ctx)
	if err != nil {
		return 0, err
	}
	if person.DeletedAt != nil {
		return 0, fmt.Errorf("person not found")
	}
	id, err := s.repository.DeletePerson(ctx, d)
	if err != nil {
		return 0, err
	}
	err = s.cache.DeletePerson(ctx)
	if err != nil {
		return 0, err
	}
	cards, err := s.repository.GetCardsByPersonID(id)
	if err != nil {
		return 0, err
	}
	for _, card := range cards {
		_, err := s.repository.DeleteCard(card.CardNumber)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}
func (s *Service) AddPerson(person *model.Person, ctx context.Context) (*model.Person, error) {
	var hashedPassword string
	hashedPassword, err := utils.HashPassword(person.Password)
	if err != nil {
		return nil, err
	}
	p, err := s.repository.GetPersonByLogin(person.Login)
	if err != nil && !errors.As(err, &model.NotFound) {
		return nil, err
	}
	if p != nil {
		return nil, fmt.Errorf("login exists")
	}
	person.Password = hashedPassword
	person, err = s.repository.AddPerson(person, ctx)
	ctx = context.WithValue(ctx, "id", person.ID)
	if err != nil {
		return nil, err
	}
	err = s.cache.SetPerson(ctx, person)
	if err != nil {
		return nil, err
	}
	return person, err
}
func (s *Service) AddPeople(people []model.Person, ctx context.Context) ([]model.Person, error) {
	people, err := s.repository.AddPeople(people, ctx)
	if err != nil {
		return nil, err
	}
	err = s.cache.SetSlice(ctx, people)
	if err != nil {
		return nil, err
	}
	return people, nil
}
func (s *Service) GetPeople(ctx context.Context) ([]model.Person, error) {
	people, err := s.repository.GetPeople(ctx)
	if err != nil {
		log.Printf("GetPeople: %v", err)
		return nil, err
	}
	return people, nil
}
func (s *Service) GetPersonByID(ctx context.Context) (*model.Person, error) {
	person, err := s.cache.GetPerson(ctx)
	if !errors.As(err, &ErrNotFoundInCache) && err != nil {
		return nil, err
	} else if errors.As(err, &ErrNotFoundInCache) {
		person, err := s.repository.GetPersonByID(ctx)
		if err != nil {
			return nil, err
		}
		var tempPerson model.Person
		tempPerson.ID = person.ID
		tempPerson.Content = person.Content
		tempPerson.Login = person.Login
		tempPerson.Password = person.CurrentPassword
		tempPerson.CreatedAt = person.CreatedAt
		tempPerson.UpdatedAt = person.UpdatedAt
		tempPerson.DeletedAt = person.DeletedAt
		err = s.cache.SetPerson(ctx, &tempPerson)
		if err != nil {
			return nil, err
		}
		return &tempPerson, err
	}
	return person, err
}
func (s *Service) UpdatePerson(person *model.Person, ctx context.Context) (*model.Person, error) {
	err := s.repository.UpdatePerson(person, ctx)
	if err != nil {
		return nil, err
	}
	_, err = s.repository.GetPersonByID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.cache.SetPerson(ctx, person)
	if err != nil {
		return nil, err
	}
	return person, nil
}
