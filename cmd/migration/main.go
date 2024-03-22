package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	goose "github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	database := os.Getenv("DATABASE_NAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database)

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("error creating connection pool: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("error setting dialect: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalf("error closing database: %v", err)
	}
}
