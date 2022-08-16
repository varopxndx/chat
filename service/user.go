package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/varopxndx/chat/model"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	queryGetUserByUsername = `SELECT * FROM users WHERE username = $1`
	queryInsertNewUser     = `INSERT INTO users VALUES(DEFAULT, $1, $2, $3)`
)

// GetUserByUsername retrieves user data
func (d DB) GetUserByUsername(ctx context.Context, userName string) (*model.User, error) {
	var user model.User

	row := d.db.QueryRowContext(ctx, queryGetUserByUsername, userName)

	var createdAt time.Time

	err := row.Scan(&user.ID, &user.UserName, &user.Password, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user %s from DB: %w", userName, err)
	}
	return &user, nil
}

// InsertUser inserts an user into DB
func (d DB) InsertUser(ctx context.Context, user model.User) error {
	// encrypt password
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 13)
	if err != nil {
		return fmt.Errorf("generating password hash: %w", err)
	}

	_, err = d.db.ExecContext(ctx, queryInsertNewUser, user.UserName, string(encryptedPass), time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("inserting user %s into DB", user.UserName))
	}
	return nil
}
