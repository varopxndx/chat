package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/varopxndx/chat/model"
)

const (
	queryGetMessages = `SELECT m.* FROM (SELECT h.message, h.created_at, users.username
						FROM chat_history as h INNER JOIN users ON h.user_id = users.id WHERE room = $1 ORDER BY h.created_at DESC LIMIT $2) m ORDER BY m.created_at ASC`
	queryInsertMessage = `INSERT INTO chat_history VALUES(DEFAULT, $1, $2, $3, $4)`
)

// GetMessages returns a list of messages
func (d DB) GetMessages(ctx context.Context, limit int, room string) ([]model.Message, error) {
	rows, err := d.db.QueryContext(ctx, queryGetMessages, room, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting messages from DB: %w", err)
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var message model.Message
		if err = rows.Scan(
			&message.Message,
			&message.CreatedAt,
			&message.User.UserName,
		); err != nil {
			return messages, fmt.Errorf("scanning row from DB error: %w", err)
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return messages, fmt.Errorf("reading rows from DB: %w", err)
	}
	return messages, nil
}

// InsertMessage store message into DB
func (d DB) InsertMessage(msg model.Message) error {
	_, err := d.db.Exec(queryInsertMessage, msg.User.ID, msg.Message, time.Now().UTC(), msg.Room)
	if err != nil {
		return fmt.Errorf("inserting message into DB: %w", err)
	}
	return nil
}
