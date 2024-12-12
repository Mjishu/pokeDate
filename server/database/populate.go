package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Animal struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Species       string `json:"species"`
	Date_of_birth string `json:"date_of_birth"`
	Sex           string `json:"sex"`
	Available     bool   `json:"available"`
	Image_src     string `json:"image_src"`
}

func PopulateDB(ctx context.Context, pool *pgxpool.Pool) {
	makeAnimals(ctx, pool)
}

func makeAnimals(ctx context.Context, pool *pgxpool.Pool) {
	sql := []string{
		` INSERT INTO animals (species, date_of_birth, sex, available, animal_type, name) VALUES (
			'Bichon Frise Poodle', '2022/10/14', 'male', true, 'dog', 'Bimbus');`,
		` INSERT INTO animals (species, date_of_birth, sex, price, available, animal_type, name) VALUES (
			'Russian Blue', '2020/09/11', 'undefined' , 9.11, false, 'cat', 'Florida');`,
		` INSERT INTO animals (species, date_of_birth, sex, price, available, animal_type, name) VALUES (
			'calico', '8008/09/11', 'female', 100.00, true, 'cat', 'Garu');`,
	}
	for _, query := range sql {

		_, err := pool.Exec(ctx, query)
		queryFail(err)
	}
}
