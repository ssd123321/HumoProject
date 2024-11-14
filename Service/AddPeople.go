package Service

import (
	"Tasks/model"
	"context"
)

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
