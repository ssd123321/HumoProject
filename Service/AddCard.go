package Service

import (
	"Tasks/model"
	"Tasks/utils"
	"context"
)

func (s *Service) AddCard(ctx context.Context, bankName string) (*model.Card, error) {
	_, err := s.GetPersonByID(ctx)
	if err != nil {
		return nil, err
	}
	var card *model.Card = &model.Card{
		PersonID:     ctx.Value("person_id").(int),
		CardNumber:   utils.GenerateCardNumber(),
		DateOfExpire: utils.GenerateExpirationDate(),
		Logotype:     utils.ChooseRandomCard(),
		Money:        0,
		BankName:     bankName,
	}
	return s.repository.AddCard(card)
}
