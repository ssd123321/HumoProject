package repository

import (
	"Tasks/model"
	"Tasks/utils"
	"gorm.io/gorm"
	"log"
)

func (r *Repository) TransferMoney(senderID int, receiverID int, sum float64, percent float64) error {
	/*
		tx := r.DB.Begin()
		err := tx.Exec("Update card SET money = money - ? Where id = ?", sum, senderID).Error
		if err != nil {
			tx.Rollback()
			log.Printf("repository - TransferMoney: %v", err)
			return err
		}
		err = tx.Exec("Update card SET money = money + ? Where id = ?", resulted, receiverID).Error
		if err != nil {
			tx.Rollback()
			log.Printf("repository - TransferMoney: %v", err)
			return err
		}
		err = tx.Exec("INSERT INTO transaction_log (sender_id, receiver_id, status, sum) VALUES (?, ?, ?, ?)", senderID, receiverID, "successfully", sum).Error
		if err != nil {
			tx.Rollback()
			log.Printf("repository - TransferMoney: %v", err)
			return err
		}
		tx.Commit()
		return nil
	*/
	result := utils.CountPercent(sum, percent)
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
		if err := tx.Exec("Update card SET money = money - ? Where id = ?", sum, senderID).Error; err != nil {
			err2 := tx.Exec("INSERT INTO transaction_log (sender_id, receiver_id, status, sum,percent, percentedsum) VALUES (?, ?, ?, ?, ?, ?)", senderID, receiverID, "unsuccessfully:"+err.Error(), sum, percent, sum-result).Error
			if err2 != nil {
				return err2
			}
			return err
		}
		if err := tx.Exec("Update card SET money = money + ? Where id = ?", sum-result, receiverID).Error; err != nil {
			err2 := tx.Exec("INSERT INTO transaction_log (sender_id, receiver_id, status, sum,percent, percentedsum) VALUES (?, ?, ?, ?, ?, ?)", senderID, receiverID, "unsuccessfully:"+err.Error(), sum, percent, sum-result).Error
			if err2 != nil {
				return err2
			}
			return err
		}
		err := tx.Transaction(func(tx2 *gorm.DB) error {
			if err := tx2.Exec("Update card SET money = money + ? Where id = ?", result, model.CorporativeCard.ID).Error; err != nil {
				err2 := tx.Exec("INSERT INTO transaction_log (sender_id, receiver_id, status, sum,percent, percentedsum) VALUES (?, ?, ?, ?, ?, ?)", senderID, receiverID, "unsuccessfully:"+err.Error(), sum, percent, sum-result).Error
				if err2 != nil {
					return err2
				}
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO transaction_log (sender_id, receiver_id, status, sum, percent, percentedsum) VALUES (?, ?, ?, ?, ?, ?)", senderID, receiverID, "successfully", sum, percent, sum-result).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO transaction_log (sender_id, receiver_id, status, sum, percent, percentedsum) VALUES (?, ?, ?, ?, ?, ?)", senderID, model.CorporativeCard.ID, "successfully", result, 0, result).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func (r *Repository) TakePercent(sum float64, percent float64) (float64, float64, error) {
	result := utils.CountPercent(sum, percent)
	tx := r.DB.Begin()
	err := tx.Exec("Update card SET money = money + ? Where id = ?", result, model.CorporativeCard.ID).Error
	if err != nil {
		tx.Rollback()
		log.Printf("repository - TakePercent: %v", err)
		return 0, 0, err
	}
	tx.Commit()
	return sum, sum - result, nil
}
