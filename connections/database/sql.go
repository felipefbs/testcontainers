package database

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/felipefbs/testcontainers/message"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewConnection(connectionString string) *Database {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("failed to connect with database", err)

		return nil
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS data_points (id SERIAL PRIMARY KEY, inserted_at timestamp DEFAULT CURRENT_TIMESTAMP, data varchar(255))")
	if err != nil {
		slog.Error("failed to create table", err)
	}

	return &Database{db: db}
}

func (db *Database) Save(ctx context.Context, data message.Message) error {
	_, err := db.db.ExecContext(ctx, "INSERT INTO data_points (sended_at, received_at, message) values ($1, $2, $3)", data.SendedAt, data.ReceivedAt, data.Message)
	if err != nil {
		slog.Error("failed to save data", err)

		return err
	}

	return nil
}
