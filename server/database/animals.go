package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Shot struct {
	Id          int
	Description string
	Name        string
}

type AnimalShot struct {
	Id          int
	Description string
	Name        string
	Next_due    time.Time
	Date_given  time.Time
}

type Animal struct {
	Id            uuid.UUID
	Name          string
	Species       string
	Date_of_birth time.Time
	Sex           string
	Price         *float32
	Available     bool
	Breed         string
	Image_src     *string
	Shots         []AnimalShot
}

type NewAnimal struct {
	Name          string          `json:"Name"`
	Species       string          `json:"Species"`
	Date_of_birth time.Time       `json:"Date_of_birth"`
	Sex           string          `json:"Sex"`
	Price         float32         `json:"Price"`
	Available     bool            `json:"Available"`
	Breed         string          `json:"Breed"`
	Image_src     string          `json:"Image_src"`
	Shots         []NewAnimalShot `json:"Shots"`
}

type NewAnimalShot struct {
	Animal_id  uuid.UUID `json:"Animal_id"`
	Shot_id    int       `json:"Shot_id"`
	Date_given time.Time `json:"Date_given"`
	Date_due   time.Time `json:"Next_due"`
}

func GetAnimalOrganization(pool *pgxpool.Pool, animal_id uuid.UUID) (uuid.UUID, error) {
	sql := `SELECT u.id FROM users u LEFT JOIN organization_animals oa ON u.id = oa.organization_id WHERE animal_id = $1`

	var organizationId uuid.UUID
	row := pool.QueryRow(context.TODO(), sql, animal_id)
	err := row.Scan(&organizationId)
	fmt.Printf("organization id is %v\n", organizationId)
	if err != nil {
		return uuid.UUID{}, err
	}
	return organizationId, nil
}

func AddUserAnimalSeen(pool *pgxpool.Pool, user_id, animal_id uuid.UUID, liked bool) error {
	sql := `INSERT INTO users_animals_seen (user_id, animal_id, liked) VALUES ($1,$2,$3)`

	_, err := pool.Exec(context.TODO(), sql, user_id, animal_id, liked)
	if err != nil {
		return err
	}
	return nil
}

func ResetSeenProgress(pool *pgxpool.Pool, user_id uuid.UUID) error {
	sql := `DELETE FROM users_animals_seen WHERE user_id = $1 AND liked = false`

	_, err := pool.Exec(context.TODO(), sql, user_id)
	if err != nil {
		return err
	}
	return nil
}

func InsertAnimal(pool *pgxpool.Pool, animal Animal) {
	sql := `
		INSERT INTO animals (name,species,date_of_birth,sex,price,available,breed) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		);
	`
	_, err := pool.Exec(context.TODO(), sql, animal.Name, animal.Species, animal.Date_of_birth, animal.Sex, animal.Price, animal.Available, animal.Breed)
	inserQueryFail(err, "inserting animal")
}

func InsertAnimalShots(pool *pgxpool.Pool, shot NewAnimalShot) error { //* work on this to insert new shot when editing.

	_, isShot := GetShot(pool, shot.Animal_id, shot.Shot_id)
	if isShot {
		_, err := pool.Exec(context.TODO(), `UPDATE animal_shots SET next_due = $1, date_given = $2 WHERE animal_id = $3 AND shots_id = $4 `, shot.Date_due, shot.Date_given, shot.Animal_id, shot.Shot_id)
		if err != nil {
			return err
		}
	}

	//* CREATE NEW
	sql := `
		INSERT INTO animal_shots(animal_id, shots_id, date_given, next_due) VALUES ($1, $2, $3, $4)
		`
	_, err := pool.Exec(context.TODO(), sql, shot.Animal_id, shot.Shot_id, shot.Date_given, shot.Date_due)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAnimalShots(pool *pgxpool.Pool, shot NewAnimalShot) error {
	sql := `
	DELETE FROM animal_shots WHERE animal_id=$1 AND shots_id=$2`

	_, err := pool.Exec(context.TODO(), sql, shot.Animal_id, shot.Shot_id)
	if err != nil {
		return err
	}
	return nil
}
func UpdateAnimalImage(pool *pgxpool.Pool, animalId uuid.UUID, imageSrc string) error {
	sql := `UPDATE animals SET image_src = $1 WHERE id = $2`

	_, err := pool.Exec(context.TODO(), sql, animalId, imageSrc)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAnimal(pool *pgxpool.Pool, animal Animal) error {
	fmt.Printf("animal recieved is %v\n", animal)
	sql := `
		UPDATE animals set name = $1, date_of_birth = $2, price = $3, available = $4 WHERE id = $5
	`

	_, err := pool.Exec(context.TODO(), sql, animal.Name, animal.Date_of_birth, animal.Price, animal.Available, animal.Id)
	if err != nil {
		fmt.Println("error!")
		return err
	}
	return nil
}

func DeleteAnimal(pool *pgxpool.Pool, id interface{}) error {
	sql := `
		DELETE FROM animals WHERE id = $1
	`
	_, err := pool.Exec(context.TODO(), sql, id)
	if err != nil {
		return err
	}
	return nil
}

func GetAnimal(pool *pgxpool.Pool, id any) (Animal, error) {
	var animal Animal
	err := pool.QueryRow(context.TODO(), "SELECT a.*, ai.image_src FROM animals AS a LEFT JOIN animal_images as ai ON a.id = ai.animal_id WHERE a.id = $1", id).Scan(
		&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
		&animal.Available, &animal.Breed, &animal.Image_src,
	)

	if err != nil {
		return Animal{}, err
	}

	SelectShots(&animal, pool)

	return animal, nil
}

func GetAnimalByName(pool *pgxpool.Pool, animal_name string) uuid.UUID {
	var id uuid.UUID
	err := pool.QueryRow(context.TODO(), "SELECT id FROM animals WHERE name = $1", animal_name).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed for animal by name: %v\n", err)
	}
	return id
}

func GetRandomAnimal(pool *pgxpool.Pool, userId uuid.UUID) (Animal, error) {
	// where a.id a
	sql := `
        SELECT a.*, ai.image_src 
        FROM animals AS a 
        LEFT JOIN animal_images AS ai ON a.id = ai.animal_id
        LEFT JOIN users_animals_seen AS uas ON a.id = uas.animal_id AND uas.user_id = $1
        WHERE uas.animal_id IS NULL
        ORDER BY RANDOM() 
        LIMIT 1
    `

	var animal Animal
	err := pool.QueryRow(context.TODO(), sql, userId).Scan(
		&animal.Id, &animal.Name, &animal.Species, &animal.Date_of_birth, &animal.Sex, &animal.Price,
		&animal.Available, &animal.Breed, &animal.Image_src,
	)
	if err != nil {
		return Animal{}, err
	}
	SelectShots(&animal, pool)
	// err = pool.QueryRow(ctx, "SELECT image_src FROM animal_images WHERE animal_id = $1;", animal.Id).Scan(&animal.Image_src)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
	// }

	return animal, nil
}

// * CREATE NEW FUNCTION THAT GETS SHOTS AND ADDS IT TO ANIMAL
func SelectShots(animal *Animal, pool *pgxpool.Pool) {
	rows, err := pool.Query(context.TODO(), "SELECT s.*,a.date_given, a.next_due FROM animal_shots AS a LEFT JOIN shots AS s ON a.shots_id = s.id WHERE animal_id = $1", animal.Id)
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

func GetAllAnimals(pool *pgxpool.Pool) []Animal {
	var animals []Animal

	rows, err := pool.Query(context.TODO(), "SELECT a.*,  ai.image_src FROM animals AS a LEFT JOIN animal_images AS ai ON a.id = ai.animal_id")
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
		SelectShots(&animal, pool)
		rowNumber++
		animals = append(animals, animal)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "row iteration of all animals failed %v\n", err)
	}
	return animals
}

func GetAllShots(pool *pgxpool.Pool) []Shot {
	var shots []Shot

	rows, err := pool.Query(context.TODO(), "SELECT * FROM shots;")
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

func GetShot(pool *pgxpool.Pool, animal_id uuid.UUID, shot_id int) (NewAnimalShot, bool) {
	var shot NewAnimalShot

	err := pool.QueryRow(context.TODO(), "SELECT animal_id, shots_id, date_given, next_due FROM animal_shots WHERE animal_id = $1 AND shots_id = $2", animal_id, shot_id).Scan(&shot.Animal_id, &shot.Shot_id, &shot.Date_given,
		&shot.Date_due)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed for get shot: %v\n", err)
		return NewAnimalShot{}, false
	}

	return shot, shot.Animal_id != uuid.UUID{}
}
