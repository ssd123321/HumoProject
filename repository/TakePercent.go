package repository

import (
	"Tasks/model"
	"Tasks/utils"
	"log"
)

func (r *Repository) TakePercent(sum float64, percent float64) (float64, error) {
	result := utils.CountPercent(sum, percent)
	tx := r.DB.Begin()
	err := tx.Exec("Update card SET money = money + ? Where id = ?", result, model.CorporativeCard.ID).Error
	if err != nil {
		tx.Rollback()
		log.Printf("repository - TakePercent: %v", err)
		return 0, err
	}
	tx.Commit()
	return sum - result, nil
}
