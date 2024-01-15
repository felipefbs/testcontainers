package database_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/felipefbs/testcontainers/connections/database"
	"github.com/felipefbs/testcontainers/message"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test_Database(t *testing.T) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithInitScripts(filepath.Join("script", "bootstrap.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}

	connString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	db := database.NewConnection(connString)
	assert.NotNil(t, db)
	msg := message.Message{
		SendedAt:   time.Now(),
		ReceivedAt: time.Now(),
		InsertedAt: time.Now(),
		Message:    "message",
	}

	err = db.Save(ctx, msg)
	assert.Nil(t, err)

	pg, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	var id int
	var inserted_at time.Time
	var data string
	err = pg.QueryRowContext(ctx, "SELECT id, inserted_at, message from data_points where message = 'message'").Scan(&id, &inserted_at, &data)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 1, id)
	assert.False(t, inserted_at.IsZero())
	assert.Equal(t, "message", data)

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()
}
