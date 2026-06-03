package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jjma22/finance-tracker/internal/data"
)

func GetTotal(uuid string) (float32, error) {

	DB.l.Info("Getting total expenses from database")
	// Runs query on database
	rows, err := DB.pool.Query(context.Background(), "SELECT price FROM expenses WHERE uuid = $1", uuid)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		DB.l.Error("Failed querying row", "error", err)
		return 0, err
	}

	//empty var for total expenses to be calculated
	var sum float32
	// Iterate through the rows from db query and scans into value to be added for total expense
	for rows.Next() {
		var n string
		err = rows.Scan(&n)
		if err != nil {
			DB.l.Error("Failed to scan rows in value", "error", err)
			return 0, err
		}
		//convert type string (from db) to float32
		i, _ := strconv.ParseFloat(n, 32)
		f := float32(i)
		sum += f
	}
	if rows.Err() != nil {
		DB.l.Error("Failed to scan rows", "error", err)
		return 0, rows.Err()
	}
	return sum, nil
}

// temp object - will remove once date columns are updated to time type on data.Expense Struct
type tempExpense struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
	// Type  string `json:"type"`
	Price      float32    `json:"price" validate:"gt=0"`
	SKU        string     `json:"sku" validate:"required,sku"`
	DateAdded  *time.Time `json:"-"`
	LastUpdate *time.Time `json:"-"`
	Uuid       string     `json:"uuid"`
}

// Fucntion to return all expenses from database
func GetExpenses(uuid string) (*data.Expenses, error) {
	DB.l.Info("Getting all expenses from database")
	// Get all ids from expense
	rows, err := DB.pool.Query(context.Background(), "SELECT id FROM expenses WHERE uuid = $1", uuid)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return nil, err
	}

	// Scan all ids into slice
	r, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (int, error) {
		var n int
		err := row.Scan(&n)
		return n, err
	})

	expenses := data.Expenses{}

	// Add expense into expenses for each id
	for _, id := range r {
		e, err := GetExpense(id, uuid)

		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)

	}

	return &expenses, nil
}

func GetExpense(id int, uuid string) (*data.Expense, error) {

	idExists, err := CheckExpenseExists(id, uuid)

	if err != nil {
		return nil, err
	}
	// Error if expense to update is not found
	if idExists == false {
		return nil, errors.New("Invalid ID")
	}

	DB.l.Info("Getting id: %d expenses from database", id)
	// Runs query on database
	row, err := DB.pool.Query(context.Background(), "SELECT * FROM expenses WHERE id = $1 AND uuid = $2", id, uuid)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return nil, err
	}

	// Parse rows in struct
	exp, err := pgx.CollectRows(row, pgx.RowToStructByName[tempExpense])
	if err != nil {
		DB.l.Error("Failed querying row", "error", err)
		return nil, err
	}

	// Prevents kernel error if last update or date added is nil
	if exp[0].DateAdded == nil {
		DB.l.Info("Setting DateAdded to nil")
		exp[0].DateAdded = &time.Time{}
	}
	if exp[0].LastUpdate == nil {
		DB.l.Info("Setting Lasttupdate to date added")
		exp[0].LastUpdate = exp[0].DateAdded
	}

	// return expense
	return &data.Expense{
		ID:         exp[0].ID,
		Name:       exp[0].Name,
		Price:      exp[0].Price,
		SKU:        exp[0].SKU,
		DateAdded:  *exp[0].DateAdded,
		LastUpdate: *exp[0].LastUpdate,
	}, nil
}

func AddExpense(e *data.Expense, uuid string) error {
	_, err := DB.pool.Exec(context.Background(), "INSERT INTO expenses (name, price, sku, dateadded,lastupdate, uuid) Values ($1, $2, $3, $4, $5, $6)",
		e.Name, e.Price, e.SKU, e.DateAdded, e.LastUpdate, uuid)

	if err != nil {
		return err
	}
	return nil

}

func DeleteExpense(id int, uuid string) (int, error) {
	ct, err := DB.pool.Exec(context.Background(), "DELETE FROM expenses WHERE id = $1 AND uuid = $2",
		id, uuid)
	if err != nil {
		return 1, err
	}

	// return 0 if no rows updated (id not found)
	if ct.RowsAffected() == 0 {
		return 0, nil
	}

	return 1, nil

}

func CheckExpenseExists(id int, uuid string) (bool, error) {
	// search fields to check if id exists
	row, err := DB.pool.Query(context.Background(), "SELECT id FROM expenses WHERE id = $1 AND uuid = $2", id, uuid)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return false, err
	}

	// Parse ids into slice
	i, err := pgx.CollectRows(row, pgx.RowTo[int])
	if err != nil {
		DB.l.Error("Failed querying row", "error", err)
		return false, err
	}

	// Return false if expense to update is not found
	if len(i) == 0 {
		return false, nil
	}
	return true, nil
}

func UpdateExpense(e *data.Expense, uuid string) error {

	id := e.ID

	idExists, err := CheckExpenseExists(id, uuid)

	if err != nil {
		return err
	}
	// Error if expense to update is not found
	if idExists == false {
		return errors.New("Invalid ID")
	}

	// Convert id into string id to be used in db querys
	sId := strconv.Itoa(id)

	// Update price of expense if value is not 0 / nil
	if e.Price != 0 {
		// Need to add validation for prices (not minus number, not string)
		lastUpdate := time.Now().Truncate(time.Second)
		DB.l.Info("Updating price for expense", "id", sId)
		u, err := DB.pool.Exec(context.Background(),
			"UPDATE expenses SET price = $1, lastupdate = $2  WHERE id = $3 AND uuid = $4",
			e.Price, lastUpdate, sId, uuid)

		if err != nil {
			return err
		}

		// return 0 if no rows updated (id not found)
		if u.RowsAffected() == 0 {
			return errors.New("Dabase update failed")
		}

	}

	// Update name of expense if value is not "" / nil
	if e.Name != "" {
		// Need to add validation for Name (not minus number, not string)
		lastUpdate := time.Now().Truncate(time.Second)
		DB.l.Info("Updating name for expense", "id", sId)
		u, err := DB.pool.Exec(context.Background(),
			"UPDATE expenses SET name = $1, lastupdate = $2  WHERE id = $3 AND uuid = $4",
			e.Name, lastUpdate, sId, uuid)

		if err != nil {
			return err
		}

		// return 0 if no rows updated (id not found)
		if u.RowsAffected() == 0 {
			return errors.New("Dabatse update failed")
		}

	}
	return nil

}
