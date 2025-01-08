package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Organization struct {
	Id       *uuid.UUID `json:"Id"`
	Name     string     `json:"Name"`
	Password string     `json:"Password"`
	Email    *string    `json:"Email"`
}

func CreateOrganization(pool *pgxpool.Pool, org Organization) error {
	sql := `
		INSERT INTO organization (name,password,email) VALUES ($1, $2, $3)
	`

	_, err := pool.Exec(context.TODO(), sql, org.Name, org.Password, org.Email)
	if err != nil {
		return err
	}
	return nil
}

func GetOrganization(pool *pgxpool.Pool, id uuid.UUID) Organization {
	sql := `
		SELECT id, name, email FROM organization WHERE id = $1
	`

	var organization Organization
	pool.QueryRow(context.TODO(), sql, id)
	return organization
}

// * this returns nil nil
func GetOrganizationByName(pool *pgxpool.Pool, name string) Organization {
	sql := `
		SELECT id, name, email, password FROM organization WHERE name = $1
	`
	var organization Organization
	pool.QueryRow(context.TODO(), sql, name).Scan(&organization.Id, &organization.Name, &organization.Email)

	fmt.Printf("stored organization is %v\n", organization)
	return organization
}

func UpdateOrganization(pool *pgxpool.Pool, org Organization) error {
	sql := `
		UPDATE organization SET name = $1, email = $2 WHERE id = $3
	`
	_, err := pool.Exec(context.TODO(), sql, org.Name, org.Email, org.Id)
	if err != nil {
		return err
	}
	return nil
}

func GetOrganizationAnimals(pool *pgxpool.Pool, id uuid.UUID) ([]Animal, error) {
	sql := `
		SELECT a.id,a.name, ai.breed FROM organization_animals oa LEFT JOIN animals a WHERE oa.organization_id = $1
	`

	rows, err := pool.Query(context.TODO(), sql, id)
	if err != nil {
		return []Animal{}, err
	}

	var animals []Animal
	for rows.Next() {
		var animal Animal
		err := rows.Scan(&animal.Id, &animal.Name, &animal.Breed)
		if err != nil {
			return []Animal{}, err
		}
		animals = append(animals, animal)
	}

	return animals, nil
}
