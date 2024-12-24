package repository

import "log"

func (r *Repository) ChangePassword(newPassword string, id int) error {
	err := r.DB.Exec("Update person SET current_password = ? Where id = ?", newPassword, id).Error
	if err != nil {
		log.Printf("ChangePassword - %v", err)
		return err
	}
	return nil
}
