package models

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/OR-Sasaki/cat-backend/config"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

func GetUser(ctx context.Context, id string) (*User, error) {
	user, err := gorm.G[User](config.DB).Where("id = ?", id).First(ctx)
	return &user, err
}

func RegisterUser(ctx context.Context, name string) (user *User, password string, err error) {
	password = uuid.New().String()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil || len(passwordHash) == 0 {
		slog.Error("failed to generate password hash", "error", err)
		return nil, "", err
	}

	u := User{
		Name:         name,
		PasswordHash: string(passwordHash),
	}

	err = gorm.G[User](config.DB).Create(ctx, &u)

	return &u, password, err
}
