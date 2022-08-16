package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/varopxndx/chat/handler/mocks/mocks"
	"github.com/varopxndx/chat/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestController_Login(t *testing.T) {
	type mockCalls struct {
		login int
		token int
	}

	tests := []struct {
		name          string
		body          model.LoginData
		err           error
		expectedCode  int
		expectedCalls mockCalls
		response      *model.User
	}{
		{
			name: "fail, bad request",
			body: model.LoginData{
				UserName: "someuser",
			},
			err:          errors.New("error: missing data"),
			expectedCode: http.StatusBadRequest,
			expectedCalls: mockCalls{
				login: 0,
				token: 0,
			},
		},
		{
			name: "fail, error calling usecase",
			body: model.LoginData{
				UserName: "someuser",
				Password: "somepassword",
				Room:     "free",
			},
			err:          errors.New("error: unable to login"),
			expectedCode: http.StatusBadRequest,
			expectedCalls: mockCalls{
				login: 1,
				token: 0,
			},
		},
		{
			name: "success",
			body: model.LoginData{
				UserName: "someuser",
				Password: "somepassword",
				Room:     "free",
			},
			err:          nil,
			expectedCode: http.StatusOK,
			expectedCalls: mockCalls{
				login: 1,
				token: 1,
			},
			response: &model.User{
				ID:       1,
				UserName: "doe",
				Password: "somepassword",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)

			ctx, _ := gin.CreateTestContext(w)

			postData, err := json.Marshal(tt.body)
			if err != nil {
				assert.Fail(t, "cannot unmarshal body: %v", err)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(postData))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			usecaseMock := &mocks.UseCase{}
			usecaseMock.On("Login", mock.AnythingOfType("*context.emptyCtx"), tt.body).Return(tt.response, tt.err)

			tokenMock := &mocks.Tokenizer{}
			tokenMock.On("GenerateToken", mock.Anything).Return("sometoken", nil)
			h := Handler{
				usecase:   usecaseMock,
				webSocket: nil,
				tokenizer: tokenMock,
			}
			h.Login(ctx)

			assert.Equal(t, tt.expectedCode, w.Code)
			usecaseMock.AssertNumberOfCalls(t, "Login", tt.expectedCalls.login)
			tokenMock.AssertNumberOfCalls(t, "GenerateToken", tt.expectedCalls.token)
		})
	}
}
