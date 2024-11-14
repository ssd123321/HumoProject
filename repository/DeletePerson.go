package repository

import (
	"context"
	"log"
	"time"
)

func (r *Repository) DeletePerson(ctx context.Context, time1 *time.Time) (int, error) {
	var id int
	err := r.DB.WithContext(ctx).Exec("UPDATE person set deleted_at = ? Where id = ? and deleted_at is null", time.Now(), ctx.Value("id").(int)).Error
	if err != nil {
		log.Printf("Failed to DeletePerson: %v", err)
		return 0, err
	}
	id = ctx.Value("id").(int)
	return id, nil
}
