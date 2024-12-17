package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// shot needs name description and id

type Shot struct {
	Id          string
	Description string
	Name        string
	Date_given  string
	Next_due    string
}

type Animal struct {
	Id            string
	Name          string
	Species       string
	Date_of_birth time.Time
	Sex           string
	Price         *float32
	Available     bool
	Breed         string
	Image_src     *string
	Shots         []Shot
}

func PopulateDB(ctx context.Context, pool *pgxpool.Pool) {
	// makeAnimals(ctx, pool)
}

func makeAnimals(ctx context.Context, pool *pgxpool.Pool) {
	sql := []string{
		` INSERT INTO animals (species, date_of_birth, sex, available, breed, name) VALUES (
			'dog', '2022/10/14', 'male', true, 'Bichon Frise Poodle', 'Bimbus');`,
		` INSERT INTO animals (species, date_of_birth, sex, price, available, breed, name) VALUES (
			'cat', '2020/09/11', 'undefined' , 9.11, false, 'Russian Blue', 'Florida');`,
		` INSERT INTO animals (species, date_of_birth, sex, price, available, breed, name) VALUES (
			'cat', '8008/09/11', 'female', 100.00, true, 'calico', 'Garu');`,
	}
	for _, query := range sql {

		_, err := pool.Exec(ctx, query)
		queryFail(err, "insert animal data")
	}
}

func makeImage(ctx context.Context, pool *pgxpool.Pool) {
	sql := []string{
		`INSERT INTO animal_images (animal_id, image_src) VALUES ('5883f423-30ee-46e3-abf8-413f1f55bdc1', './images/dog.webp');`,
	}

	for _, query := range sql {
		_, err := pool.Exec(ctx, query)
		queryFail(err, "insert image data")
	}
}
