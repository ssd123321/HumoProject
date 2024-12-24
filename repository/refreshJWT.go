package repository

import "fmt"

func (r *Repository) AddRefresh(id int, refreshJWT string) error {
	fmt.Println(123)
	tx := r.DB.Begin()
	result := tx.Table("refresh_token").Where("person_id = ?", id).Update("token", refreshJWT)
	if result.RowsAffected == 0 {
		err := tx.Exec("INSERT INTO refresh_token (person_id, token) VALUES (?, ?)", id, refreshJWT).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
func (r *Repository) GetRefreshByPersonID(id int) (string, error) {
	var refreshJWT string
	err := r.DB.Raw("Select refresh_token.token From refresh_token Where person_id = ?", id).First(&refreshJWT).Error
	if err != nil {
		return "", err
	}
	return refreshJWT, nil
}
