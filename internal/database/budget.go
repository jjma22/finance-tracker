package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jjma22/finance-tracker/internal/data"
)

// Set new budget
func SetBudget(b int, uuid string) error {
	// Set current data and time
	dateAdded := time.Now().Truncate(time.Second)
	lastUpdate := time.Now().Truncate(time.Second)

	// Insert into database
	_, err := DB.pool.Exec(context.Background(), "INSERT INTO budget (budget,uuid,dateadded,lastupdate) Values ($1, $2, $3, $4",
		b, uuid, dateAdded, lastUpdate)
	if err != nil {
		return err
	}

	DB.l.Info("Successfully added budget into database")
	return nil
}

// Get budget from id in path
func GetBudget(id int, uuid string) (*data.Budget, error) {

	// Run query on database to return budget
	row, err := DB.pool.Query(context.Background(), "select budget from budget where id = $1 AND uuid = $2", id, uuid)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return nil, err
	}

	// Read budget rows into slice b
	b, err := pgx.CollectRows(row, pgx.RowTo[int])
	if err != nil {
		DB.l.Error("Failed querying row", "error", err)
		return nil, err
	}

	if len(b) == 0 {
		return nil, errors.New("Invalid id, no results returned")
	}

	if len(b) > 1 {
		return nil, errors.New("multiple values returned for id")
	}
	// return budget
	return &data.Budget{
		Budget: b[0],
	}, nil

}

// Update budget with id from URL path
func UpdateBudget(id int, b int, uuid string) error {
	// Need to add check to see if id exists
	// Run query on database to update budget
	r, err := DB.pool.Exec(context.Background(), "UPDATE budget SET budget = $1 WHERE id = $2 and uuid = $3", b, id, uuid)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return err
	}

	// return 0 if no rows updated (id not found)
	if r.RowsAffected() == 0 {
		return errors.New("Database update failed")
	}

	return nil

}
