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
		&animal.Available, &animal.Breed,
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
		&animal.Available, &animal.Breed,
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

/* func GetOrganizationAnimal() []Animal
^ Get organization id from body
^ search in organization_animals for all animals from X organization
^ Get all the animals data from that organization
^ return the slice of animals
*/

func GetAllAnimals() []Animal {
	ctx, pool := createConnection()
	var animals []Animal

	//* how to store the animal info into animals?``
	rows, err := pool.Query(ctx, "SELECT * FROM animals")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var animal Animal
		err := rows.Scan(&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
			&animal.Available, &animal.Breed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan row: %v\n", err)
			continue
		}
		animals = append(animals, animal)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "row iteration failed %v\n", err)
	}
	return animals
}
