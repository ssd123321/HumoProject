package repository

import (
	"Tasks/model"
	"Tasks/time"
	"context"
)

func (r *Repository) UpdatePerson(person *model.Person, ctx context.Context) error {
	currentTime, err := time.GetCurrentDate()
	if err != nil {
		return err
	}
	err = r.DB.WithContext(ctx).Exec("UPDATE Person Set content = ?, updated_at = ?  Where id = ? and person.deleted_at is null returning content::text", person.Content, currentTime, ctx.Value("id")).Error
	if err != nil {
		return err
	}
	return nil
}
