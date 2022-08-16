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

func TestController_SignUp(t *testing.T) {
	tests := []struct {
		name         string
		body         model.User
		err          error
		expectedCode int
	}{
		{
			name: "fail",
			body: model.User{
				ID:       1,
				UserName: "doe",
				Password: "somepassword",
			},
			expectedCode: http.StatusBadRequest,
			err:          errors.New("unable to signup"),
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

			req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(postData))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			usecaseMock := &mocks.UseCase{}
			usecaseMock.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), tt.body).Return(tt.err)

			h := Handler{
				usecase:   usecaseMock,
				webSocket: nil,
				tokenizer: nil,
			}
			h.SignUp(ctx)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
