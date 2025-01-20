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
