package database

import (
	"fmt"
	"os"
)

func GetLocations() Location {
	ctx, pool := createConnection()
	var location Location
	err := pool.QueryRow(ctx, "SELECT * FROM locations").Scan(
		&location.Id,
		&location.Name,
		&location.Location_type,
		&location.Parent_id,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	}

	fmt.Println(location)
	return location
}

func GetAnimals(id string) Animal {
	ctx, pool := createConnection()
	var animal Animal
	err := pool.QueryRow(ctx, "SELECT * from animals where id = $1", id).Scan(
		&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
		&animal.Available, &animal.Animal_type,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	}

	err = pool.QueryRow(ctx, "SELECT image_src FROM animal_images WHERE animal_id = $1;", animal.Id).Scan(&animal.Image_src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	}

	return animal
}

func GetRandomAnimal() Animal {
	ctx, pool := createConnection()
	var animal Animal
	err := pool.QueryRow(ctx, "SELECT * FROM animals ORDER BY RANDOM() LIMIT 1").Scan(
		&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
		&animal.Available, &animal.Animal_type,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	}
	err = pool.QueryRow(ctx, "SELECT image_src FROM animal_images WHERE animal_id = $1;", animal.Id).Scan(&animal.Image_src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	}

	return animal
}
