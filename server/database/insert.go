package database

import (
	"fmt"
	"os"
	"time"
)

type NewAnimal struct {
	Name          string              `json:"name"`
	Species       string              `json:"species"`
	Date_of_birth time.Time           `json:"date_of_birth"`
	Sex           string              `json:"sex"`
	Price         float32             `json:"price"`
	Available     bool                `json:"available"`
	Breed         string              `json:"breed"`
	Shots         []NewShotFromClient `json:"shots"`
}

type NewAnimalShot struct {
	Animal_id  string    `json:"animal_id"`
	Shot_id    string    `json:"shot_id"`
	Date_given time.Time `json:"date_given"`
	Date_due   time.Time `json:"date_due"`
}

type NewShotFromClient struct {
	Shot_id    string    `json:"id"`
	Date_given time.Time `json:"date_given"`
	Date_due   time.Time `json:"date_due"`
}

type UpdateAnimalStruct struct {
	Id            string              `json:"id"`
	Name          string              `json:"name"`
	Date_of_birth time.Time           `json:"date_of_birth"`
	Price         float32             `json:"price"`
	Available     bool                `json:"available"`
	Image_src     string              `json:"image_src"`
	Shots         []NewShotFromClient `json:"shots"`
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

func InsertAnimalShots(shot NewAnimalShot) { //! inserting Date_due gives incorrect date! fix this.!
	//* if animal alrady as shot THEN update the shot with new info
	ctx, pool := createConnection()

	_, isShot := GetShot(shot.Animal_id)
	if isShot {
		fmt.Println("is shot is true")
		_, err := pool.Exec(ctx, `UPDATE animal_shots SET next_due = $1, date_give = $2 WHERE animal_id = $2 AND shots_id = $3 `, shot.Date_due, shot.Date_given, shot.Animal_id, shot.Shot_id)
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

func UpdateAnimal(animal UpdateAnimalStruct) { //! Add ability to add new shots
	sql := `
		UPDATE animals  SET name = $1, date_of_birth = $2, price = $3, available = $4  WHERE id = $5
	`
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, sql, animal.Name, animal.Date_of_birth, animal.Price, animal.Available, animal.Id)
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
