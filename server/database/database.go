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

func CreateConnection() (context.Context, *pgxpool.Pool) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, GetItemFromENV("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}
	return ctx, dbpool
}

func Database(pool *pgxpool.Pool) {
	// defer pool.Close()

	// callSchemas(context.TODO(), pool)
	// PopulateDB(context.TODO(), pool)
}

func GetItemFromENV(key string) string {
	err := godotenv.Load("/.env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			err = godotenv.Load(".env")
			if err != nil {
				fmt.Printf("could not load .env")
			}
		}
	}

	dbURL := os.Getenv(key)
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in .env")
	}
	return dbURL
}
