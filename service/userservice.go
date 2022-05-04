package service

import (
	"context"
	"fmt"

	"github.com/Msksgm/go-lightweight-realworld.git/model"
	"github.com/Msksgm/go-lightweight-realworld.git/repository"
)

type UserService struct {
	userRepository repository.UserRepositorier
}

func NewUserService(userRepository repository.UserRepositorier) (*UserService, error) {
	return &UserService{userRepository: userRepository}, nil
}

func (us *UserService) RegisterUser(ctx context.Context, email string, userName string, password string) error {
	user, err := us.userRepository.FindUserByUserName(ctx, userName)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("userName %v is already used", userName)
	}
	user, err = us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("email %v is already used", email)
	}
	user = &model.User{Email: email, UserName: userName, PasswordHash: password}
	err = us.userRepository.SaveUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
