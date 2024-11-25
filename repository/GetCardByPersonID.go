package repository

import (
	"Tasks/model"
	"context"
)

func (r *Repository) GetCardByPersonID(ctx context.Context) (*model.Card, error) {
	var card model.Card
	err := r.DB.WithContext(ctx).First(&card, ctx.Value("id").(int)).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}
