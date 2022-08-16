package usecase

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/varopxndx/chat/model"

	"github.com/pkg/errors"
)

const (
	stock            = "stock"
	csvNumFields     = 8
	stockBotResponse = "%s quote is $%s per share"
)

var commands = map[string]bool{
	stock: true,
}

// GetMessages retrieves last 50 messages from DB
func (u *UseCase) GetMessages(ctx context.Context, limit int, room string) ([]model.Message, error) {
	return u.db.GetMessages(ctx, limit, room)
}

// ProcessMessage handle the message
func (u *UseCase) ProcessMessage(message string) ([]byte, error) {
	var msg model.Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return []byte(message), errors.Wrap(err, "unmarshalling json message")
	}

	// if message is command, validate command
	if strings.HasPrefix(msg.Message, "/") {
		return u.runCommand(message, msg)
	}

	// insert message into DB
	err := u.db.InsertMessage(msg)
	if err != nil {
		return []byte(message), errors.Wrap(err, "inserting message into database")
	}

	return []byte(message), nil
}

func (u *UseCase) runCommand(message string, msg model.Message) ([]byte, error) {
	res := model.Message{
		ID: 0,
		User: model.User{
			ID:       0,
			UserName: "chat-bot",
		},
		CreatedAt: time.Now().UTC(),
	}
	c, v, err := parseCommand(msg.Message)
	if err != nil {
		res.Message = err.Error()
		b, _ := json.Marshal(res)
		return b, errors.Wrap(err, message)
	}

	if !isKnownCommand(c) {
		res.Message = fmt.Sprintf("non supported command: %s", c)
		b, _ := json.Marshal(res)
		return b, errors.New(res.Message)
	}

	price, err := u.callStooqAPI(v)
	if err != nil {
		res.Message = "error calling stooq api..."
		b, _ := json.Marshal(res)
		return b, errors.Wrap(err, "calling stooq API")
	}

	res.Message = fmt.Sprintf(stockBotResponse, strings.ToUpper(v), price)
	err = u.publishStockMessage(res)
	if err != nil {
		res.Message = "error executing command"
		b, _ := json.Marshal(res)
		return b, errors.Wrap(err, "publishing message to rabbit broker")
	}
	return nil, nil
}

func (u *UseCase) publishStockMessage(msg model.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return u.rabbit.Publish(u.rabbit.GetQueueName(), b)
}

func isKnownCommand(command string) bool {
	return commands[command]
}

func parseCommand(message string) (command, value string, err error) {
	v := strings.Split(message[1:], "=")
	if len(v) != 2 {
		return "", "", fmt.Errorf("invalid format")
	}
	command = strings.ToLower(strings.TrimSpace(v[0]))
	value = strings.ToLower(strings.TrimSpace(v[1]))
	return
}

func (u *UseCase) callStooqAPI(stockCode string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(u.stooqUrlString, stockCode))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = csvNumFields
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}
	return records[1][6], nil
}
