package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjma22/finance-tracker.git/internal/data"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "expenses"
)

type db struct {
	l    *slog.Logger
	pool *pgxpool.Pool
}

var DB = db{}

func newDb(l *slog.Logger) error {
	url := "postgresql://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		slog.Error("Could not connect to database -", "Error", err)
		return err
	}
	config.MaxConns = 4
	config.MinConns = 0
	config.MaxConnIdleTime = time.Minute * 5
	config.HealthCheckPeriod = time.Minute * 1

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("Could create database pool-", "Error", err)
		return err
	}
	DB.l = l
	DB.pool = pool
	return nil
}

// Func to try connection to db multiple times. Panics after tring 3 times
func InitDb(l *slog.Logger) {
	i := 0
	for i < 2 {
		err := newDb(l)

		if err == nil {
			DB.l.Info("Successfully established database connection")
			break
		}

		slog.Error("Error to connect to db, trying again")
		time.Sleep(5 * time.Second)

		i++

		if i == 2 {
			slog.Error("Timed out trying to connect to databse, shutting down")
			os.Exit(1)
		}
	}
}

func GetTotal() (float32, error) {

	DB.l.Info("Getting total expenses from database")
	// Runs query on database
	rows, err := DB.pool.Query(context.Background(), "select price from expenses")

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
}

// Fucntion to return all expenses from database
func GetExpenses() (*data.Expenses, error) {
	DB.l.Info("Getting all expenses from database")
	// Get all ids from expense
	rows, err := DB.pool.Query(context.Background(), "select id from expenses")
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
		e, err := GetExpense(id)

		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)

	}

	return &expenses, nil
}

func GetExpense(id int) (*data.Expense, error) {

	DB.l.Info("Getting (id - needs adding) expenses from database")
	// Runs query on database
	row, err := DB.pool.Query(context.Background(), "select * from expenses where id = $1", id)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return nil, err
	}

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

	return &data.Expense{
		ID:         exp[0].ID,
		Name:       exp[0].Name,
		Price:      exp[0].Price,
		SKU:        exp[0].SKU,
		DateAdded:  *exp[0].DateAdded,
		LastUpdate: *exp[0].LastUpdate,
	}, nil
}

func AddExpense(e *data.Expense) error {
	_, err := DB.pool.Exec(context.Background(), "INSERT INTO expenses (name, price, sku, dateadded,lastupdate) Values ($1, $2, $3, $4, $5)",
		e.Name, e.Price, e.SKU, e.DateAdded, e.LastUpdate)

	if err != nil {
		return err
	}
	return nil

}

func DeleteExpense(id int) (int, error) {
	ct, err := DB.pool.Exec(context.Background(), "DELETE FROM expenses WHERE id = $1",
		id)
	if err != nil {
		return 1, err
	}

	// return 0 if no rows updated (id not found)
	if ct.RowsAffected() == 0 {
		return 0, nil
	}

	return 1, nil

}

func UpdateExpense(e *data.Expense) error {
	// search fields to check if id exists
	row, err := DB.pool.Query(context.Background(), "select id from expenses where id = $1", e.ID)
	if err != nil {
		DB.l.Error("Failed querying database", "error", err)
		return err
	}

	i, err := pgx.CollectRows(row, pgx.RowTo[int])
	if err != nil {
		DB.l.Error("Failed querying row", "error", err)
		return err
	}

	if len(i) == 0 {
		return errors.New("Invalid ID")
	}

	sId := strconv.Itoa(i[0])

	if e.Price != 0 {
		// Need to add validation for prices (not minus number, not string)
		lastUpdate := time.Now().Truncate(time.Second)
		DB.l.Info("Updating price for expense", "id", sId)
		u, err := DB.pool.Exec(context.Background(),
			"UPDATE expenses SET price = $1, lastupdate = $2  WHERE id = $3",
			e.Price, lastUpdate, sId)

		if err != nil {
			return err
		}

		// return 0 if no rows updated (id not found)
		if u.RowsAffected() == 0 {
			return errors.New("Dabatse update failed")
		}

	}

	if e.Name != "" {
		// Need to add validation for Name (not minus number, not string)
		lastUpdate := time.Now().Truncate(time.Second)
		DB.l.Info("Updating name for expense", "id", sId)
		u, err := DB.pool.Exec(context.Background(),
			"UPDATE expenses SET name = $1, lastupdate = $2  WHERE id = $3",
			e.Name, lastUpdate, sId)

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
