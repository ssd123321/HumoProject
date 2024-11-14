package Service

import (
	"Tasks/time"
	"context"
	"fmt"
	"log"
)

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
	return id, nil
}
