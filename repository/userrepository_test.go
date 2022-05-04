package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
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
			user: &model.User{Email: "test@example.com", UserName: "test-user", PasswordHash: "password"},
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO users").
					WithArgs("test@example.com", "test-user", "password").
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
				return mock{db, m}
			}(),
			want:   nil,
			hasErr: false,
		},
		{
			name: "fail",
			user: &model.User{Email: "test@example.com", UserName: "test-user", PasswordHash: "password"},
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO users").
					WithArgs("test@example.com", "test-user", "password").
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

func TestFindUserByUserName(t *testing.T) {
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
	}
	tests := []struct {
		name     string
		userName string
		mock     mock
		want     *model.User
		hasErr   bool
		wantErr  error
	}{
		{
			name:     "found",
			userName: "test-user",
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, username, password FROM users WHERE username = $1`)).
					WithArgs("test-user").
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
						AddRow("0", "test@example.com", "test-user", "password"),
					)
				m.ExpectCommit()
				return mock{db, m}
			}(),
			want:    &model.User{Email: "test@example.com", UserName: "test-user", PasswordHash: "password"},
			hasErr:  false,
			wantErr: nil,
		},
		{
			name:     "not found",
			userName: "test-user",
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, username, password FROM users WHERE username = $1`)).
					WithArgs("test-user").
					WillReturnRows(sqlmock.NewRows([]string{}))
				m.ExpectCommit()
				return mock{db, m}
			}(),
			want:    nil,
			hasErr:  false,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.mock.db
			userRepository, err := NewUserRepository(db)
			if err != nil {
				t.Error(err)
			}
			got, err := userRepository.FindUserByUserName(context.Background(), tt.userName)
			if (err != nil) != tt.hasErr {
				t.Errorf("err type: %v, expect err type: %v", reflect.TypeOf(err), reflect.TypeOf(tt.wantErr))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindUserByEmail(t *testing.T) {
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
	}
	tests := []struct {
		name    string
		email   string
		mock    mock
		want    *model.User
		hasErr  bool
		wantErr error
	}{
		{
			name:  "found",
			email: "test@example.com",
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, username, password FROM users WHERE email = $1`)).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
						AddRow("0", "test@example.com", "test-user", "password"),
					)
				m.ExpectCommit()
				return mock{db, m}
			}(),
			want:    &model.User{Email: "test@example.com", UserName: "test-user", PasswordHash: "password"},
			hasErr:  false,
			wantErr: nil,
		},
		{
			name:  "not found",
			email: "test@example.com",
			mock: func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				m.ExpectBegin()
				m.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, username, password FROM users WHERE email = $1`)).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{}))
				m.ExpectCommit()
				return mock{db, m}
			}(),
			want:    nil,
			hasErr:  false,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.mock.db
			userRepository, err := NewUserRepository(db)
			if err != nil {
				t.Error(err)
			}
			got, err := userRepository.FindUserByEmail(context.Background(), tt.email)
			if (err != nil) != tt.hasErr {
				t.Errorf("err type: %v, expect err type: %v", reflect.TypeOf(err), reflect.TypeOf(tt.wantErr))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
