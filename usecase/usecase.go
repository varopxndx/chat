package usecase

import (
	"context"

	"github.com/varopxndx/chat/model"
)

//go:generate mockery --name Service --filename ./mocks/service.go --outpkg mocks
type Service interface {
	GetUserByUsername(ctx context.Context, userName string) (*model.User, error)
	InsertUser(ctx context.Context, user model.User) error
	GetMessages(ctx context.Context, limit int, room string) ([]model.Message, error)
	InsertMessage(msg model.Message) error
}

//go:generate mockery --name Broker --filename ./mocks/broker.go --outpkg mocks
type Broker interface {
	GetQueueName() string
	Publish(key string, msg []byte) error
}

// UseCase struct
type UseCase struct {
	db             Service
	rabbit         Broker
	stooqUrlString string
}

// New creates the usecase
func New(db Service, rabbit Broker, stooqUrlString string) *UseCase {
	return &UseCase{db: db, rabbit: rabbit, stooqUrlString: stooqUrlString}
}
