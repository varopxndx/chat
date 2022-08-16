package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/varopxndx/chat/model"
	"github.com/varopxndx/chat/usecase/mocks/mocks"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_Login(t *testing.T) {
	login := model.LoginData{
		UserName: "doe",
		Password: "somePassword",
		Room:     "free",
	}
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("somePassword"), 13)
	user := &model.User{
		UserName: "doe",
		Password: string(encryptedPassword),
	}

	tests := []struct {
		name             string
		userData         model.LoginData
		expectedResponse *model.User
		err              error
		wantErr          bool
		errResponse      error
	}{
		{
			name:             "success",
			userData:         login,
			expectedResponse: user,
		},
		{
			name:             "fail, unable to retrieve user",
			userData:         login,
			expectedResponse: user,
			wantErr:          true,
			err:              errors.New("unable to retrieve user"),
		},
		{
			name:             "fail, user not registered",
			userData:         login,
			expectedResponse: nil,
			wantErr:          true,
			errResponse:      errors.New("username doe not registered"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repoMock := &mocks.Service{}
			repoMock.On("GetUserByUsername", ctx, tt.userData.UserName).Return(tt.expectedResponse, tt.err)

			u := &UseCase{
				db: repoMock,
			}
			_, err := u.Login(ctx, tt.userData)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repoMock.AssertExpectations(t)
		})
	}
}
