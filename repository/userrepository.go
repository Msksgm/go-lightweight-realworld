package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Msksgm/go-lightweight-realworld.git/model"
)

type UserRepositorier interface {
	SaveUser(context.Context, *model.User) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

type SaveUserQueryError struct {
	User    *model.User
	Message string
	Err     error
}

func (err *SaveUserQueryError) Error() string {
	return err.Message
}

func (ur *UserRepository) SaveUser(ctx context.Context, user *model.User) (err error) {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(query, user.Email, user.UserName, user.PasswordHash)
	if err != nil {
		return &SaveUserQueryError{User: user, Message: fmt.Sprintf("userrepository.SaveUser err: %s", err), Err: err}
	}
	return nil
}
