package mock

import (
	"context"

	"github.com/Msksgm/go-lightweight-realworld.git/model"
)

type UserRepositoryStub struct {
	SaveUserFn           func(ctx context.Context, user *model.User) error
	FindUserByUserNameFn func(ctx context.Context, userName string) (*model.User, error)
	FindUserByEmailFn    func(ctx context.Context, email string) (*model.User, error)
}

func (urs *UserRepositoryStub) SaveUser(ctx context.Context, user *model.User) error {
	return urs.SaveUserFn(ctx, user)
}

func (urs *UserRepositoryStub) FindUserByUserName(ctx context.Context, userName string) (*model.User, error) {
	return urs.FindUserByUserNameFn(ctx, userName)
}

func (urs *UserRepositoryStub) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return urs.FindUserByEmailFn(ctx, email)
}
