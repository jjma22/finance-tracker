package database

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
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
	host string
	port int
}

func NewDb(l *slog.Logger) (*db) {
	return &db{
		l: l,
		host : host,
		port: port,

	}
}

func (d*db) ConnectDb() (error) {
	url := "postgresql://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		d.l.Error("Could not connect to database -", "Error", err)
		return err
	}
	d.l.Info("Successfully established connection to database")
	defer conn.Close(context.Background())
	return nil

}