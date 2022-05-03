package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint
	Email        string
	UserName     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uint, email string, userName string, password string, createdAt time.Time, updatedAt time.Time) (*User, error) {
	hasBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{ID: id, Email: email, UserName: userName, PasswordHash: string(hasBytes), CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
