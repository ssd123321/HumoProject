package repository

import "Tasks/model"

func (r *Repository) AddCard(c *model.Card) (*model.Card, error) {
	err := r.DB.Table("card").Create(&c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}
