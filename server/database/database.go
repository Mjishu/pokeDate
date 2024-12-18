package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Location struct {
	Id            int
	Name          string
	Location_type string
	Parent_id     *int
}

func createConnection() (context.Context, *pgxpool.Pool) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, GetDatabaseURL())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}
	return ctx, dbpool
}

func Database() {
	ctx, dbpool := createConnection()
	defer dbpool.Close()

	// callSchemas(ctx, dbpool)
	PopulateDB(ctx, dbpool)
	// getLocations(ctx, dbpool)
}

func GetDatabaseURL() string {
	err := godotenv.Load("../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal("ERror loading .env file")
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in .env")
	}
	return dbURL
}
