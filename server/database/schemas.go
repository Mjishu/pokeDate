package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func callSchemas(ctx context.Context, pool *pgxpool.Pool) {
	createLocations(ctx, pool)
	createUsers(ctx, pool)
	createOrganization(ctx, pool)
	createShots(ctx, pool)
	createAnimals(ctx, pool)
}

// * DONE SO FAR: locations, users, organizations, shots

// need to create the enum as well
func createLocations(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS locations (
			id SERIAL PRIMARY KEY, 
			name VARCHAR(100), 
			location_type location_type NOT NULL, 
			parent_id INT REFERENCES locations(id) ON DELETE SET NULL
		);
	`

	_, err := pool.Exec(ctx, sql)
	queryFail(err)
}

func createAnimals(ctx context.Context, pool *pgxpool.Pool) {
	// add likes dislikes ?location?
	sql := `
		CREATE TABLE IF NOT EXISTS animals (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(30) NOT NULL,
			species VARCHAR(100) NOT NULL, 
			date_of_birth DATE NOT NULL,
			sex sex_enum NOT NULL,
			price FLOAT,
			available BOOLEAN NOT NULL,
			animal_type VARCHAR(50) NOT NULL
		);
	`

	_, err := pool.Exec(ctx, sql)
	queryFail(err)
}

func createShots(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS shots (
			id SERIAL PRIMARY KEY,
			name VARCHAR(150) NOT NULL,
			description TEXT 
		)
	`

	_, err := pool.Exec(ctx, sql)
	queryFail(err)
}

func createUsers(ctx context.Context, pool *pgxpool.Pool) {
	// add interested in tags?
	sql := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(40) NOT NULL,
			email VARCHAR(100) ,
			password VARCHAR(50) NOT NULL,
			date_of_birth DATE,
			country_id INT REFERENCES locations(id) ON DELETE SET NULL,
			state_id INT REFERENCES locations(id) ON DELETE SET NULL,
			city_id INT REFERENCES locations(id) ON DELETE SET NULL,
			profile_picture_src TEXT
		);
	`
	_, err := pool.Exec(ctx, sql)
	queryFail(err)
}

func createOrganization(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS organization (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(50) NOT NULL,
			password VARCHAR(100) NOT NULL,
			country_id INT REFERENCES locations(id) ON DELETE SET NULL,
			state_id INT REFERENCES locations(id) ON DELETE SET NULL,
			city_id INT REFERENCES locations(id) ON DELETE SET NULL,
			website_url text
		)
	`
	_, err := pool.Exec(ctx, sql)
	queryFail(err)
}

func queryFail(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Table 'users' created successfully")
}

// ! TODO: ADD animal_shots_join, organization_animals, animal_images

// ! less todo: add messages, conversation, conversation_member
