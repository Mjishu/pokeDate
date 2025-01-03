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
	createAnimalImages(ctx, pool)
	createAnimalShots(ctx, pool)
	createOrganizationAnimals(ctx, pool)
	createUserAnimals(ctx, pool)
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
	queryFail(err, "locations")
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
			price FLOAT DEFAULT 0.00,
			available BOOLEAN NOT NULL,
			breed VARCHAR(50) NOT NULL
		);
	`

	_, err := pool.Exec(ctx, sql)
	queryFail(err, "animals")
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
	queryFail(err, "shots")
}

func createUsers(ctx context.Context, pool *pgxpool.Pool) {
	// add interested in tags?
	sql := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(40) NOT NULL,
			email VARCHAR(100) ,
			password text NOT NULL,
			date_of_birth DATE,
			country_id INT REFERENCES locations(id) ON DELETE SET NULL,
			state_id INT REFERENCES locations(id) ON DELETE SET NULL,
			city_id INT REFERENCES locations(id) ON DELETE SET NULL,
			profile_picture_src TEXT
			created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
			updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
		);
	`
	_, err := pool.Exec(ctx, sql)
	queryFail(err, "users")
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
	queryFail(err, "organization")
}

func createAnimalImages(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS animal_images (
			animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
			image_src TEXT NOT NULL,
			priority INT NOT NULL, 
			image_public_id VARCHAR(120) NOT NULL,
			CONSTRAINT unique_priority_per_animal UNIQUE (animal_id, priority)
		);
	`

	_, err := pool.Exec(ctx, sql)
	queryFail(err, "animal Images")
}

func createAnimalShots(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS animal_shots (
			animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
			shots_id BIGINT REFERENCES shots(id) ON DELETE CASCADE NOT NULL,
			date_given DATE,
			next_due DATE,
			PRIMARY KEY (animal_id, shots_id)
		);
	`
	_, err := pool.Exec(ctx, sql)
	queryFail(err, "animal shots")
}

func createOrganizationAnimals(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS organization_animals (
			organization_id UUID REFERENCES organization(id) ON DELETE CASCADE NOT NULL,
			animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
			PRIMARY KEY (organization_id, animal_id)
		);
	`
	_, err := pool.Exec(ctx, sql)
	queryFail(err, "organization animals")
}

func createUserAnimals(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS user_animals (
			user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
			liked BOOLEAN NOT NULL,
			PRIMARY KEY (user_id, animal_id)
		);
	`
	_, err := pool.Exec(ctx, sql)
	queryFail(err, "user animals")
}

func createRefreshTokens(ctx context.Context, pool *pgxpool.Pool) {
	sql := `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			token TEXT PRIMARY KEY NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
			updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL,
			revoked_at TIMESTAMPTZ DEFAULT NULL
		)
	`
	_, err := pool.Exec(ctx, sql)
	queryFail(err, "refresh tokens")
}

// ! todo: add messages, conversation, conversation_member

func queryFail(err error, tableName string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Table '%s' created successfully\n", tableName)
}
