package Service

import (
	"Tasks/model"
	"context"
)

func (s *Service) AddCard(c *model.Card, ctx context.Context) (*model.Card, error) {
	_, err := s.repository.GetCardByPersonID(ctx)
	if err != nil {
		return nil, err
	}
	c, err = s.repository.AddCard(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
