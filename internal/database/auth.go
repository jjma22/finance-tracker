package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jjma22/finance-tracker/internal/auth"
)

func GetUser(username string) (*auth.User, error) {
	DB.l.Info("Getting user from database")

	row, err := DB.pool.Query(context.Background(), "select id,username,password from users where username = $1", username)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return nil, err
	}

	tempUser, err := pgx.CollectRows(row, pgx.RowToStructByName[auth.User])
	if err != nil {
		DB.l.Error("Failed querying row", "error", err)
		return nil, err
	}

	if len(tempUser) == 0 {
		return nil, errors.New("No user found")
	}

	if len(tempUser) > 1 {
		return nil, errors.New("multiple users found, username is not unique")
	}

	return &tempUser[0], nil
}
