package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/varopxndx/chat/model"
	"github.com/varopxndx/chat/usecase/mocks/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_GetMessages(t *testing.T) {
	ctx := context.Background()

	user := model.User{
		UserName: "doe",
		Password: "somePassword",
	}

	tests := []struct {
		name     string
		messages []model.Message
		limit    int
		err      error
		wantErr  bool
	}{
		{
			name: "success",
			messages: []model.Message{
				{
					ID:        1,
					User:      user,
					Message:   "Some message",
					CreatedAt: time.Time{},
				},
			},
			limit: 50,
		},
		{
			name:     "success, no messages found",
			messages: nil,
			limit:    50,
		},
		{
			name:     "fail",
			messages: nil,
			limit:    50,
			err:      errors.New("Database error..."),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.Service{}
			repoMock.On("GetMessages", ctx, tt.limit, "free").Return(tt.messages, tt.err)

			u := UseCase{
				db: &repoMock,
			}

			got, err := u.GetMessages(ctx, tt.limit, "free")

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.messages, got)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestUseCase_ProcessMessage(t *testing.T) {
	user := model.User{
		UserName: "doe",
		Password: "somePassword",
	}
	created := time.Now().Format(time.RFC3339)
	parsedTime, _ := time.Parse(time.RFC3339, created)
	message := model.Message{
		User:      user,
		Message:   "some message",
		CreatedAt: parsedTime,
	}
	commandMessage := model.Message{
		User:      user,
		Message:   "/stock=someStockCode",
		CreatedAt: parsedTime,
	}

	type mockCalls struct {
		insert   int
		publish  int
		getStock int
	}

	tests := []struct {
		name        string
		message     model.Message
		err         error
		wantErr     bool
		expectCalls mockCalls
	}{
		{
			name:    "success",
			message: message,
			expectCalls: mockCalls{
				insert:   1,
				publish:  0,
				getStock: 0,
			},
		},
		{
			name:    "success, command message",
			message: commandMessage,
			expectCalls: mockCalls{
				insert:   0,
				publish:  1,
				getStock: 1,
			},
		},
		{
			name:    "fail, database error",
			message: message,
			wantErr: true,
			err:     errors.New("database error..."),
			expectCalls: mockCalls{
				insert:   1,
				publish:  0,
				getStock: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mocks.Service)
			repoMock.On("InsertMessage", message).Return(tt.err)

			brokerMock := new(mocks.Broker)
			brokerMock.On("Publish", "testQueue", mock.Anything).Return(nil)
			brokerMock.On("GetQueueName").Return("testQueue")

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`Symbol,Date,Time,Open,High,Low,Close,Volume
								test,2021-11-10,01:00:00,0,0,0,100.50,0`))
			}))
			defer server.Close()

			u := UseCase{
				db:             repoMock,
				rabbit:         brokerMock,
				stooqUrlString: server.URL + "/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv",
			}

			b, _ := json.Marshal(tt.message)
			_, err := u.ProcessMessage(string(b))

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repoMock.AssertNumberOfCalls(t, "InsertMessage", tt.expectCalls.insert)
			brokerMock.AssertNumberOfCalls(t, "Publish", tt.expectCalls.publish)
			brokerMock.AssertNumberOfCalls(t, "GetQueueName", tt.expectCalls.getStock)
		})
	}
}
