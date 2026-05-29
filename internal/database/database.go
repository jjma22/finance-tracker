package database

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	env_config "github.com/jjma22/finance-tracker/internal/config"
)

type db struct {
	l    *slog.Logger
	pool *pgxpool.Pool
}

var DB = db{}

func newDb(l *slog.Logger, db *env_config.Database) error {

	// should make ssl mode a flag
	url := "postgresql://" + db.DB_user + ":" + db.DB_password + "@" + db.DB_host + ":" + db.DB_port + "/" + db.DB_name + "?sslmode=disable"
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
func InitDb(l *slog.Logger, db *env_config.Database) {
	i := 0
	for i < 2 {
		err := newDb(l, db)

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
