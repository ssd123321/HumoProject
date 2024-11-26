package repository

import (
	"Tasks/model"
	"log"
)

func (r *Repository) TransferMoney(sender *model.Card, receiver *model.Card, sum float64) error {
	tx := r.DB.Begin()
	err := tx.Exec("Update card SET money = money - ? Where card_number = ?", sender.Money, sender.CardNumber).Error
	if err != nil {
		tx.Rollback()
		log.Printf("repository - TransferMoney: %v", err)
		return err
	}
	err = tx.Exec("Update card SET money = money + ? Where card_number = ?", receiver.Money, receiver.CardNumber).Error
	if err != nil {
		tx.Rollback()
		log.Printf("repository - TransferMoney: %v", err)
		return err
	}
	err = tx.Exec("INSERT INTO transaction_log (sender_id, receiver_id, status) VALUES (?, ?, ?)", sender.ID, receiver.ID, "successfully").Error
	if err != nil {
		tx.Rollback()
		log.Printf("repository - TransferMoney: %v", err)
		return err
	}
	tx.Commit()
	return nil
}
