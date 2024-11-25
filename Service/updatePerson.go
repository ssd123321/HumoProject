package Service

import (
	"Tasks/model"
	"context"
)

func (s *Service) UpdatePerson(person *model.Person, ctx context.Context) (*model.Person, error) {
	err := s.repository.UpdatePerson(person, ctx)
	if err != nil {
		return nil, err
	}
	person, err = s.repository.GetPersonByID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.cache.SetPerson(ctx, person)
	if err != nil {
		return nil, err
	}
	return person, nil
}
