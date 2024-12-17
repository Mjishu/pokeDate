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
	Breed         string    `json:"breed"`
}

type UpdateAnimalStruct struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Date_of_birth time.Time `json:"date_of_birth"`
	Price         float32   `json:"price"`
	Available     bool      `json:"available"`
}

func InsertAnimal(animal NewAnimal) { //! fix
	sql := `
		INSERT INTO animals (name,species,date_of_birth,sex,price,available,breed) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		);
	`
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, sql, animal.Name, animal.Species, animal.Date_of_birth, animal.Sex, animal.Price, animal.Available, animal.Breed)
	inserQueryFail(err, "inserting animal")
}

func UpdateAnimal(animal UpdateAnimalStruct) {
	sql := `
		UPDATE animals  SET name = $1, date_of_birth = $2, price = $3, available = $4  WHERE id = $5
	`
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, sql, animal.Name, animal.Date_of_birth, animal.Price, animal.Available, animal.Id)
	inserQueryFail(err, "Updating Animal")
}

func inserQueryFail(err error, name string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("command  '%s' created successfully\n", name)
}
