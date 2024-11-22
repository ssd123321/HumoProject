package Service

import (
	"Tasks/model"
	"context"
	"errors"
)

func (s *Service) GetPersonByID(ctx context.Context) (*model.Person, error) {
	person, err := s.cache.GetPerson(ctx)
	if !errors.As(err, &ErrNotFoundInCache) && err != nil {
		return nil, err
	} else if errors.As(err, &ErrNotFoundInCache) {
		person, err := s.repository.GetPersonByID(ctx)
		if err != nil {
			return nil, err
		}
		err = s.cache.SetPerson(ctx, person)
		if err != nil {
			return nil, err
		}
		return person, err
	}
	return person, err
}
