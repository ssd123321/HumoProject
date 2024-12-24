package Service

import (
	"context"
	"fmt"
)

func (s *Service) ExecuteTransaction(senderNumber int, receiverNumber int, sum float64, ctx context.Context) error {
	sen, err := s.repository.GetCardByNumber(senderNumber)
	if err != nil {
		return err
	} else if sen.Money < sum {
		return fmt.Errorf("you don't have enough money")
	} else if sen.ID != ctx.Value("id").(int) {
		return fmt.Errorf("not persmissible")
	}
	rec, err := s.repository.GetCardByNumber(receiverNumber)
	if err != nil {
		return err
	}

	err = s.repository.TransferMoney(sen.ID, rec.ID, sum, 10)
	if err != nil {
		return err
	}
	return nil
}
