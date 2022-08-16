package handler

import (
	"context"

	"github.com/varopxndx/chat/model"
)

//go:generate mockery --name UseCase --filename ./mocks/usecase.go --outpkg mocks
// UseCase methods
type UseCase interface {
	Login(ctx context.Context, data model.LoginData) (*model.User, error)
	SignUp(ctx context.Context, data model.User) error
	GetMessages(ctx context.Context, limit int, room string) ([]model.Message, error)
}

//go:generate mockery --name Tokenizer --filename ./mocks/token.go --outpkg mocks
// Tokenizer methods
type Tokenizer interface {
	GenerateToken(user model.User) (string, error)
}

// Handler struct
type Handler struct {
	usecase   UseCase
	webSocket *Socket
	tokenizer Tokenizer
}

// New creates a new Handler
func New(u UseCase, ws *Socket, t Tokenizer) Handler {
	return Handler{usecase: u, webSocket: ws, tokenizer: t}
}
