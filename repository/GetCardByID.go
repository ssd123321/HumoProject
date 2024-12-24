package repository

import (
	"Tasks/model"
	"log"
	"time"
)

func (r *Repository) GetCardByNumber(number int) (*model.Card, error) {
	var card model.Card
	err := r.DB.Raw("Select * FROM card where card_number = ? and deleted_at is null", number).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}
func (r *Repository) DeleteCard(number int) (int, error) {
	err := r.DB.Exec("Update card SET deleted_at = ? Where card_number = ?", time.Now(), number).Error
	if err != nil {
		log.Printf("repository - DeleteCard: %v", err)
		return 0, err
	}
	return number, nil
}
func (r *Repository) GetCardsByPersonID(id int) ([]model.Card, error) {
	var cards []model.Card
	err := r.DB.Raw("Select * From card where person_id = ? and deleted_at is null", id).Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}
func (r *Repository) AddCard(c *model.Card) (*model.Card, error) {
	log.Println(c)
	err := r.DB.Table("card").Create(&c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}
