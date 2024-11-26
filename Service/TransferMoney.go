package Service

import "Tasks/model"

func (s *Service) ExecuteTransaction(sender *model.Card, receiver *model.Card, sum float64) error {
	sum, err := s.repository.TakePercent(sum, 2)
	if err != nil {
		return err
	}
	err = s.repository.TransferMoney(sender, receiver, sum)
	if err != nil {
		return err
	}
	return nil
}
