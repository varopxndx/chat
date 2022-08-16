package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/varopxndx/chat/model"

	"golang.org/x/crypto/bcrypt"
)

// Login validates user login information
func (u *UseCase) Login(ctx context.Context, data model.LoginData) (*model.User, error) {
	user, err := u.db.GetUserByUsername(ctx, data.UserName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("username %s not registered", data.UserName)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}
	user.Password = ""
	return user, nil
}
