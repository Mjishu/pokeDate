package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Organization struct {
	Id       uuid.UUID `json:"Id"`
	Name     string    `json:"Name"`
	Password string    `json:"Password"`
	Email    *string   `json:"Email"`
}

func CreateOrganization(pool *pgxpool.Pool, org Organization) (uuid.UUID, error) {
	sql := `
		INSERT INTO users (username,password,email, is_organization) VALUES ($1, $2, $3, true) RETURNING id
	`

	//* checks if org with same name exists
	var orgExists string
	row := pool.QueryRow(context.TODO(), "SELECT * FROM users WHERE username = $1 AND is_organization = true", org.Name)
	err := row.Scan(&orgExists)
	if err == nil {
		return uuid.UUID{}, errors.New("an organization with that name already exists")
	}

	var orgId uuid.UUID
	pool.QueryRow(context.TODO(), sql, org.Name, org.Password, org.Email).Scan(&orgId)
	if (orgId == uuid.UUID{}) {
		return uuid.UUID{}, errors.New("new organization id is empty")
	}
	return orgId, nil
}

func GetOrganization(pool *pgxpool.Pool, id uuid.UUID) Organization {
	sql := `
		select id, username, email from users where id = $1 AND is_organization = true
	`

	var organization Organization
	pool.QueryRow(context.TODO(), sql, id).Scan(&organization.Id, &organization.Name, &organization.Email)
	return organization
}

// * this returns nil nil
func GetOrganizationByName(pool *pgxpool.Pool, name string) Organization { //* should be fixed, got 4 items instead of 3 so maybe
	sql := `
		SELECT id, username FROM users WHERE username = $1 AND is_organization = true
	`
	var organization Organization
	pool.QueryRow(context.TODO(), sql, name).Scan(&organization.Id, &organization.Name)

	fmt.Printf("stored organization is %v\n", organization) // this gives big null error?
	return organization
}

func UpdateOrganization(pool *pgxpool.Pool, org Organization) error {
	sql := `
		UPDATE users SET username = $1, email = $2 WHERE id = $3 AND is_organization = true
	`
	_, err := pool.Exec(context.TODO(), sql, org.Name, org.Email, org.Id)
	if err != nil {
		return err
	}
	return nil
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
