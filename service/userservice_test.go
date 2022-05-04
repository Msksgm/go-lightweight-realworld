package service

import (
	"context"
	"testing"

	"github.com/msksgm/go-lightweight-realworld/mock"
	"github.com/msksgm/go-lightweight-realworld/model"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name               string
		email              string
		userName           string
		password           string
		saveUser           func(context.Context, *model.User) error
		findUserByUserName func(context.Context, string) (*model.User, error)
		findUserByEmail    func(context.Context, string) (*model.User, error)
		want               error
		errMsg             string
	}{{
		"success",
		"test@exmple.com",
		"test-user",
		"password",
		func(ctx context.Context, u *model.User) error {
			return nil
		},
		func(_ context.Context, s string) (*model.User, error) {
			return nil, nil
		},
		func(_ context.Context, s string) (*model.User, error) {
			return nil, nil
		},
		nil,
		"",
	}}

	userService := UserService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService.userRepository = &mock.UserRepositoryStub{SaveUserFn: tt.saveUser, FindUserByUserNameFn: tt.findUserByUserName, FindUserByEmailFn: tt.findUserByEmail}
			err := userService.RegisterUser(context.Background(), tt.email, tt.userName, tt.password)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != tt.errMsg {
				t.Errorf("Expected error `%s`, got `%s`", tt.errMsg, errMsg)
			}
		})
	}
}
