package manager

import (
	"context"
	"fmt"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/auth"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
)

type UserManager struct {
}

func (mgr *UserManager) UserRegister(ctx context.Context, user types.UserRegisterRequest) error {
	log.Debugln(ctx, "UserManager.UserRegister is called")
	if err := auth.PasswordSecurityCheck(user.Password); err != nil {
		log.Errorf(ctx, "Password security check error:%v", err)
		return err
	}
	hashedPass := auth.HashPassword(user.Password)
	if user.DisplayName == "" {
		user.DisplayName = user.UserName
	}
	u := &model.User{
		UserName:    user.UserName,
		UserPass:    hashedPass,
		DisplayName: user.DisplayName,
		Party:       user.Party,
		IsRoot:      false,
	}
	if err := model.AddUser(u); err != nil {
		log.Errorf(ctx, "Register user %s error:%v\n", user.UserName, err)
		return fmt.Errorf("Register user %s error:%w", user.UserName, err)
	}
	return nil
}

func (mgr *UserManager) UserList(ctx context.Context) ([]*types.User, error) {
	retUsers := make([]*types.User, 0)
	users, err := model.ListUsers()
	if err != nil {
		return retUsers, err
	}

	for _, u := range users {
		if !u.IsRoot { // only list member user
			apiU := &types.User{
				UserName: u.UserName,
				Party:    u.Party,
			}
			retUsers = append(retUsers, apiU)
		}
	}
	return retUsers, nil
}

func (mgr *UserManager) UserDel(ctx context.Context, userName string) error {
	user, err := model.GetUserByUserName(userName)
	if err != nil {
		return err
	}
	return model.DeleteUser(user)
}
