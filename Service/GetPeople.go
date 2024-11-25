package Service

import (
	"Tasks/model"
	"context"
	"log"
)

func (s *Service) GetPeople(ctx context.Context) ([]model.Person, error) {
	people, err := s.repository.GetPeople(ctx)
	if err != nil {
		log.Printf("GetPeople: %v", err)
		return nil, err
	}
	return people, nil
}
