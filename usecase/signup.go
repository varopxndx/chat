package usecase

import (
	"context"

	"github.com/varopxndx/chat/model"
)

// SignUp inserts a new user into DB
func (u *UseCase) SignUp(ctx context.Context, data model.User) error {
	return u.db.InsertUser(ctx, data)
}
