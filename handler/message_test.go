package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/varopxndx/chat/handler/mocks/mocks"
	"github.com/varopxndx/chat/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestController_Message(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		response     []model.Message
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			response: []model.Message{
				{
					ID: 1,
					User: model.User{
						ID:       1,
						UserName: "doe",
						Password: "somepassword",
					},
					Message:   "Hi there",
					CreatedAt: time.Now(),
				},
			},
		},
		{
			name:         "fail, unable to get messages",
			expectedCode: http.StatusInternalServerError,
			response:     nil,
			err:          errors.New("unable to get messages"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)

			ctx, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("POST", "/message", nil)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req
			ctx.Request.URL, _ = url.Parse("?room=free")

			usecaseMock := &mocks.UseCase{}
			usecaseMock.On("GetMessages", mock.AnythingOfType("*context.emptyCtx"), 50, "free").Return(tt.response, tt.err)

			h := Handler{
				usecase:   usecaseMock,
				webSocket: nil,
				tokenizer: nil,
			}
			h.Message(ctx)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
