package database

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx"
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

func (*db) newDb(l *slog.Logger) (*db) {
	return &db{
		l: l,
		host : host,
		port: port,

	}
}

func (d*db) connectDb() (error) {
	conString, err := pgx.ParseConfig((host + ":" + strconv.Itoa(port) ))
	if err != nil {
		d.l.Error("Could not parse db connection string")
		return err
	}
	_ , err = pgx.Connect(context.Background(), conString)
	if err != nil {
		d.l.Error("Could not connect to database")
		return err
	}

}