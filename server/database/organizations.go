package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Organization struct {
	Id              uuid.UUID `json:"Id"`
	Name            string    `json:"Name"`
	Password        string    `json:"Password"`
	Email           *string   `json:"Email"`
	Profile_picture string
}

func GetOrganizationAnimals(pool *pgxpool.Pool, id uuid.UUID) ([]Animal, error) {
	//! ISSUE WITH THIS WHERE STATEMENT
	sql := ` 
		SELECT a.id,a.name, a.breed, a.species, a.date_of_birth, a.sex, a.price, a.available FROM organization_animals oa LEFT JOIN animals a ON oa.animal_id = a.id WHERE oa.organization_id = $1
	`

	rows, err := pool.Query(context.TODO(), sql, id)
	if err != nil {
		return []Animal{}, err
	}

	var animals []Animal
	for rows.Next() {
		var animal Animal
		err := rows.Scan(
			&animal.Id,
			&animal.Name,
			&animal.Breed,
			&animal.Species,
			&animal.Date_of_birth,
			&animal.Sex,
			&animal.Price,
			&animal.Available,
		)
		if err != nil {
			return []Animal{}, err
		}
		animals = append(animals, animal)
	}

	return animals, nil
}

func CreateNewAnimal(pool *pgxpool.Pool, animal Animal) (uuid.UUID, error) {
	sql := `
		INSERT INTO animals(name,species,date_of_birth,sex,price,available,breed) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id
	`

	var newAnimalId uuid.UUID
	pool.QueryRow(context.TODO(), sql, animal.Name, animal.Species, animal.Date_of_birth, animal.Sex, animal.Price, animal.Available, animal.Breed).Scan(&newAnimalId)
	if (newAnimalId == uuid.UUID{}) {
		return uuid.UUID{}, errors.New("error querying rows")
	}

	return newAnimalId, nil
}

func CreateOrganizationAnimal(pool *pgxpool.Pool, orgId, animalId uuid.UUID) (bool, error) {
	sql := `
		INSERT INTO organization_animals (organization_id, animal_id) VALUES ($1, $2)
	`

	_, err := pool.Exec(context.TODO(), sql, orgId, animalId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddAnimalImage(pool *pgxpool.Pool, imageURL string, animalId uuid.UUID, priority int) error {
	sql := `
		INSERT INTO animal_images (animal_id, image_src, priority) VALUES ($1,$2,$3)
	`

	_, err := pool.Exec(context.TODO(), sql, animalId, imageURL, priority)
	if err != nil {
		return err
	}

	return nil
}
