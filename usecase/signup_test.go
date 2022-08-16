package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/varopxndx/chat/model"
	"github.com/varopxndx/chat/usecase/mocks/mocks"

	"github.com/stretchr/testify/assert"
)

func TestUseCase_SignUp(t *testing.T) {
	user := model.User{
		UserName: "doe",
		Password: "somePassword",
	}

	tests := []struct {
		name     string
		userData model.User
		err      error
		wantErr  bool
	}{
		{
			name:     "success",
			userData: user,
		},
		{
			name:     "fail",
			userData: user,
			err:      errors.New("Error"),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repoMock := &mocks.Service{}
			repoMock.On("InsertUser", ctx, tt.userData).Return(tt.err)

			u := &UseCase{
				db: repoMock,
			}

			err := u.SignUp(ctx, tt.userData)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repoMock.AssertExpectations(t)
		})
	}
}
