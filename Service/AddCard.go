package Service

import (
	"Tasks/model"
	"Tasks/utils"
	"context"
)

func (s *Service) AddCard(ctx context.Context, bankName string) (*model.Card, error) {
	_, err := s.repository.GetPersonByID(ctx)
	if err != nil {
		return nil, err
	}
	var card *model.Card = &model.Card{
		PersonID:     ctx.Value("id").(int),
		CardNumber:   utils.GenerateCardNumber(),
		DateOfExpire: utils.GenerateExpirationDate(),
		Logotype:     utils.ChooseRandomCard(),
		Money:        0,
		BankName:     bankName,
	}
	return s.repository.AddCard(card)
}
func (s *Service) DeleteCard(number int) (int, error) {
	_, err := s.repository.GetCardByNumber(number)
	if err != nil {
		return 0, err
	}
	return s.repository.DeleteCard(number)
}
func (s *Service) GetCardByNumber(number int) (*model.Card, error) {
	c, err := s.repository.GetCardByNumber(number)
	if err != nil {
		return nil, err
	}
	return c, nil
}
