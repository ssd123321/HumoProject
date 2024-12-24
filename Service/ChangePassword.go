package Service

import (
	"Tasks/utils"
	"context"
	"fmt"
)

func (s *Service) ChangePassword(newPassword string, curPassword string, ctx context.Context) error {
	p, err := s.repository.GetPersonByID(ctx)
	if err != nil {
		return err
	}
	b := utils.CheckPasswordHash(curPassword, p.CurrentPassword)
	if !b {
		return fmt.Errorf("ChangePassword - Password is incorrect")
	}
	if newPassword == curPassword {
		return fmt.Errorf("password can't be the same")
	}
	fmt.Println(p.Passwords)
	if len(p.Passwords) != 0 {
		if utils.CheckPasswordHash(newPassword, p.Passwords[0]) || utils.CheckPasswordHash(newPassword, p.Passwords[1]) || utils.CheckPasswordHash(newPassword, p.Passwords[2]) {
			return fmt.Errorf("you have used this password before")
		}
	}
	newPassword, err = utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	err = s.repository.ChangePassword(newPassword, ctx.Value("id").(int))
	if err != nil {
		return err
	}
	return nil
}
