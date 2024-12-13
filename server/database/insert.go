package database

import (
	"fmt"
	"os"
	"time"
)

type NewAnimal struct {
	Name          string    `json:"name"`
	Species       string    `json:"species"`
	Date_of_birth time.Time `json:"date_of_birth"`
	Sex           string    `json:"sex"`
	Price         float32   `json:"price"`
	Available     bool      `json:"available"`
	Animal_type   string    `json:"animal_type"`
}

func InsertAnimal(animal NewAnimal) {
	sql := `
		INSERT INTO animals (name,species,date_of_birth,sex,price,available,animal_type) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		);
	`
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, sql, animal.Name, animal.Species, animal.Date_of_birth, animal.Sex, animal.Price, animal.Available, animal.Animal_type)
	inserQueryFail(err, "inserting animal")
}

func inserQueryFail(err error, name string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("command  '%s' created successfully\n", name)
}
