package Service

import (
	"Tasks/model"
	"context"
)

func (s *Service) GetPeople(ctx context.Context) ([]model.Person, error) {
	people, err := s.repository.GetPeople(ctx)
	if err != nil {
		return nil, err
	}
	return people, nil
}
