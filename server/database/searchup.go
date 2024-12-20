package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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

func GetAnimal(id any) Animal {
	ctx, pool := createConnection()
	var animal Animal
	err := pool.QueryRow(ctx, "SELECT a.*, ai.image_src FROM animals AS a LEFT JOIN animal_images as ai ON a.id = ai.animal_id WHERE a.id = $1", id).Scan(
		&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
		&animal.Available, &animal.Breed, &animal.Image_src,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	}

	SelectShots(&animal, ctx, pool)

	return animal
}

func GetAnimalByName(animal_name string) string {
	ctx, pool := createConnection()
	var id string
	err := pool.QueryRow(ctx, "SELECT id FROM animals WHERE name = $1", animal_name).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed for animal by name: %v\n", err)
	}
	return id
}

func GetRandomAnimal() Animal {
	ctx, pool := createConnection()
	var animal Animal
	err := pool.QueryRow(ctx, "SELECT a.*, ai.image_src FROM  animals AS a LEFT JOIN  animal_images AS ai ON  a.id = ai.animal_id ORDER BY  RANDOM() LIMIT 1;").Scan(
		&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
		&animal.Available, &animal.Breed, &animal.Image_src,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	}
	SelectShots(&animal, ctx, pool)
	// err = pool.QueryRow(ctx, "SELECT image_src FROM animal_images WHERE animal_id = $1;", animal.Id).Scan(&animal.Image_src)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	// }

	return animal
}

// * CREATE NEW FUNCTION THAT GETS SHOTS AND ADDS IT TO ANIMAL
func SelectShots(animal *Animal, ctx context.Context, pool *pgxpool.Pool) {
	rows, err := pool.Query(ctx, "SELECT s.*,a.date_given, a.next_due FROM animal_shots AS a LEFT JOIN shots AS s ON a.shots_id = s.id WHERE animal_id = $1", animal.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error selecting animals shots %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var shot AnimalShot
		err := rows.Scan(&shot.Id, &shot.Name, &shot.Description, &shot.Date_given, &shot.Next_due)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Scaning row for shots didn't work: %v\n", err)
			fmt.Fprintf(os.Stderr, "shot was: %v\n", shot)
			return
		}
		fmt.Printf("Your shot is; %v\n", shot)
		animal.Shots = append(animal.Shots, shot)
	}

	if rows.Err() != nil {
		fmt.Fprintf(os.Stderr, "Rows iteration of shots failed: %v\n", err)
		return
	}
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

	rows, err := pool.Query(ctx, "SELECT a.*,  ai.image_src FROM animals AS a LEFT JOIN animal_images AS ai ON a.id = ai.animal_id")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	defer rows.Close()

	rowNumber := 0
	for rows.Next() {
		var animal Animal
		err := rows.Scan(&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
			&animal.Available, &animal.Breed, &animal.Image_src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan row: %v\n", err)
			continue
		}
		SelectShots(&animal, ctx, pool)
		rowNumber++
		animals = append(animals, animal)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "row iteration of all animals failed %v\n", err)
	}
	return animals
}

func GetAllShots() []Shot {
	ctx, pool := createConnection()
	var shots []Shot

	rows, err := pool.Query(ctx, "SELECT * FROM shots;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Failed: %v\n", err)
	}
	defer rows.Close()

	rowNumber := 0
	for rows.Next() {
		var shot Shot
		err := rows.Scan(&shot.Id, &shot.Name, &shot.Description)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan row : %v\n", err)
			continue
		}
		rowNumber++
		shots = append(shots, shot)
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "row iteration of all shots failed. %v\n", err)
	}
	return shots
}

func GetShot(animal_id string) (NewAnimalShot, bool) {
	ctx, pool := createConnection()
	var shot NewAnimalShot

	err := pool.QueryRow(ctx, "SELECT aniimal_id, shots_id, date_given, next_due FROM animal_shots WHERE animal_id = $1", animal_id).Scan(&shot.Animal_id, &shot.Shot_id, &shot.Date_given,
		&shot.Date_due)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed for get shot: %v\n", err)
		return NewAnimalShot{}, false
	}
	return shot, true
}
