package repository

import (
	"Tasks/model"
	"context"
)

func (r *Repository) UpdatePerson(person *model.Person, ctx context.Context) error {
	err := r.DB.WithContext(ctx).Exec("UPDATE Person Set content = ?  Where id = ? and person.deleted_at is null returning content::text", person.Content, ctx.Value("id")).Error
	if err != nil {
		return err
	}
	return nil
}
