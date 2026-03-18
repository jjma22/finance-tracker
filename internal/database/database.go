package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "expenses"
)


type db struct {
	l *slog.Logger
	pool *pgxpool.Pool
}

 var DB = db{}

func newDb(l *slog.Logger) (error) {
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

//Func to try connection to db multiple times. Panics after tring 3 times
func InitDb(l *slog.Logger) () {
		i := 0
			for i < 2 {
		err := newDb(l)
		
		if err == nil {
			DB.l.Info("Successfully established database connection")
			break;
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

	DB.l.Info("Running database query")
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
		fmt.Println(n)
		if err != nil {
			DB.l.Error("Failed to scan rows in value", "error", err)
			return 0, err
		}
		//convert type string (from db) to float32
		i,_ := strconv.ParseFloat(n, 32)
		f := float32(i)
		sum += f
	}
	if rows.Err() != nil {
		DB.l.Error("Failed to scan rows", "error", err)
		return 0, rows.Err()
	}
	return sum, nil
}