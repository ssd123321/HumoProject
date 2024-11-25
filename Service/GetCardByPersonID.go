package Service

import (
	"Tasks/model"
	"context"
)

func (s *Service) GetCardByPersonID(ctx context.Context) (*model.Card, error) {
	c, err := s.repository.GetCardByPersonID(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}
