package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
