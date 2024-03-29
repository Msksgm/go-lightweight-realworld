package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/msksgm/go-lightweight-realworld/model"
)

type UserRepositorier interface {
	SaveUser(context.Context, *model.User) error
	FindUserByUserName(context.Context, string) (*model.User, error)
	FindUserByEmail(context.Context, string) (*model.User, error)
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
	_, err = tx.ExecContext(ctx, query, user.Email, user.UserName, user.PasswordHash)
	if err != nil {
		return &SaveUserQueryError{User: user, Message: fmt.Sprintf("userrepository.SaveUser err: %s", err), Err: err}
	}
	return nil
}

func (ur *UserRepository) FindUserByUserName(ctx context.Context, userName string) (user *model.User, err error) {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
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
		SELECT id, email, username, password FROM users WHERE username = $1
	`
	rows, err := tx.QueryContext(ctx, query, userName)
	if err != nil {
		return nil, err
	}

	var u model.User
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.UserName, &u.PasswordHash)
		if err != nil {
			return nil, err
		}
		user = &u
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) FindUserByEmail(ctx context.Context, email string) (user *model.User, err error) {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
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
		SELECT id, email, username, password FROM users WHERE email = $1
	`
	rows, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u model.User
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.UserName, &u.PasswordHash)
		if err != nil {
			return nil, err
		}
		user = &u
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}
