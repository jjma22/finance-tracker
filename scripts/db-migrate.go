package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	pool *pgxpool.Pool
}

var DB = db{}

func newDb(h string, pr string, u string, pw string, d string) error {
	url := "postgresql://" + u + ":" + pw + "@" + h + ":" + pr + "/" + d + "?sslmode=disable"
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
	DB.pool = pool
	return nil
}

// Run migrations on test database

func main() {

	host := os.Getenv("ACTIONS_DB_host")
	port := os.Getenv("ACTIONS_DB_port")
	user := os.Getenv("ACTIONS_DB_user")
	pw := os.Getenv("ACTIONS_DB_pw")
	db := os.Getenv("ACTIONS_DB_name")

	newDb(host, port, user, pw, "postgres")

	// Create expenses table
	_, err := DB.pool.Exec(context.Background(), "CREATE DATABASE expenses WITH  OWNER = postgres")
	if err != nil {
		slog.Error("Error creating expenses database", "error", err)
		os.Exit(1)
	}

	newDb(host, port, user, pw, db)

	_, err = DB.pool.Exec(context.Background(), "CREATE TABLE expenses  (id SERIAL PRIMARY KEY, name VARCHAR(255),	price NUMERIC, sku VARCHAR(255), dateadded timestamptz, lastupdate timestamptz)")
	if err != nil {
		slog.Error("Error creating expenses table", "error", err)
		os.Exit(1)
	}

	_, err = DB.pool.Exec(context.Background(), "CREATE TABLE budget  (id SERIAL PRIMARY KEY, budget NUMERIC, dateadded timestamptz, lastupdate timestamptz)")
	if err != nil {
		slog.Error("Error creating budget table", "error", err)
		os.Exit(1)
	}

}
