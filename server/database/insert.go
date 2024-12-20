package database

import (
	"fmt"
	"os"
	"time"
)

type NewAnimal struct {
	Name          string              `json:"Name"`
	Species       string              `json:"Species"`
	Date_of_birth time.Time           `json:"Date_of_birth"`
	Sex           string              `json:"Sex"`
	Price         float32             `json:"Price"`
	Available     bool                `json:"Available"`
	Breed         string              `json:"Breed"`
	Shots         []NewShotFromClient `json:"Shots"`
}

type NewAnimalShot struct {
	Animal_id  string    `json:"Animal_id"`
	Shot_id    int       `json:"Shot_id"`
	Date_given time.Time `json:"Date_given"`
	Date_due   time.Time `json:"Next_due"`
}

type NewShotFromClient struct {
	Shot_id    int       `json:"Id"`
	Date_given time.Time `json:"Date_given"`
	Date_due   time.Time `json:"Next_due"`
}

type UpdateAnimalStruct struct {
	Id            string              `json:"Id"`
	Name          string              `json:"Name"`
	Date_of_birth time.Time           `json:"Date_of_birth"`
	Price         float32             `json:"Price"`
	Available     bool                `json:"Available"`
	Image_src     string              `json:"Image_src"`
	Shots         []NewShotFromClient `json:"Shots"`
}

func InsertAnimal(animal NewAnimal) {
	sql := `
		INSERT INTO animals (name,species,date_of_birth,sex,price,available,breed) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		);
	`
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, sql, animal.Name, animal.Species, animal.Date_of_birth, animal.Sex, animal.Price, animal.Available, animal.Breed)
	inserQueryFail(err, "inserting animal")
}

func InsertAnimalShots(shot NewAnimalShot) {

	ctx, pool := createConnection()

	//! This isn't working properly to check if shot exists, i create new shot and it goes to the isShot if statement
	_, isShot := GetShot(shot.Animal_id, shot.Shot_id)
	if isShot {
		fmt.Println("is shot is true")
		_, err := pool.Exec(ctx, `UPDATE animal_shots SET next_due = $1, date_given = $2 WHERE animal_id = $3 AND shots_id = $4 `, shot.Date_due, shot.Date_given, shot.Animal_id, shot.Shot_id)
		inserQueryFail(err, "Updating shot")
		return
	}

	//* CREATE NEW
	fmt.Println("isShot was false")
	sql := `
		INSERT INTO animal_shots(animal_id, shots_id, date_given, next_due) VALUES ($1, $2, $3, $4)
		`
	_, err := pool.Exec(ctx, sql, shot.Animal_id, shot.Shot_id, shot.Date_given, shot.Date_due)
	inserQueryFail(err, "Inserting shot")
}

func UpdateAnimal(animal UpdateAnimalStruct) {
	sql := `
		UPDATE animals SET name = $1, date_of_birth = $2, price = $3, available = $4 WHERE id = $5
	`
	fmt.Printf("the updated animal is %v\n and the animal_id = %v\n", animal, animal.Id)
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, sql, animal.Name, animal.Date_of_birth, animal.Price, animal.Available, animal.Id) //? Why is this giving an error?
	inserQueryFail(err, "Updating Animal")
}

func inserQueryFail(err error, name string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed of %s: %v\n", name, err)
	}
	fmt.Printf("command  '%s' created successfully\n", name)
}

func DeleteAnimal(id interface{}) {
	sql := `
		DELETE FROM animals WHERE id = $1
	`
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, sql, id)
	inserQueryFail(err, "deleting animal")
}
