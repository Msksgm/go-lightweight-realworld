package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Msksgm/go-lightweight-realworld.git/model"
)

func TestSaveUser(t *testing.T) {
	var saveUserQueryError *SaveUserQueryError
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
	}
	tests := []struct {
		name   string
		user   *model.User
		mock   mock
		want   error
		hasErr bool
	}{
		{
			name: "success",
			user: &model.User{Email: "test@examle.com", UserName: "test-user", PasswordHash: "password"},
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO users").
					WithArgs("test@examle.com", "test-user", "password").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
				return mock{db, m}
			}(),
			want:   nil,
			hasErr: false,
		},
		{
			name: "fail",
			user: &model.User{Email: "test@examle.com", UserName: "test-user", PasswordHash: "password"},
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO users").
					WithArgs("test@examle.com", "test-user", "password").
					WillReturnError(saveUserQueryError)
				m.ExpectRollback()
				return mock{db, m}
			}(),
			want:   saveUserQueryError,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.mock.db
			userRepository, err := NewUserRepository(db)
			if err != nil {
				t.Error(err)
			}
			got := userRepository.SaveUser(context.Background(), tt.user)
			if (got != tt.want) != tt.hasErr {
				t.Errorf("got %s, want %s", got, tt.want)
			}
			if errors.As(got, &tt.want) != tt.hasErr {
				t.Errorf("err type: %v, expect err type: %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}
