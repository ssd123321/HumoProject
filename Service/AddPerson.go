package Service

import (
	"Tasks/model"
	"context"
)

func (s *Service) AddPerson(person *model.Person, ctx context.Context) (*model.Person, error) {
	person, err := s.repository.AddPerson(person, ctx)
	ctx = context.WithValue(ctx, person.ID, person.ID)
	if err != nil {
		return nil, err
	}
	err = s.cache.SetPerson(ctx, person)
	if err != nil {
		return nil, err
	}
	return person, err
}
