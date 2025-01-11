package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateRefreshToken(pool *pgxpool.Pool, token string, id uuid.UUID) (bool, error) {

	sql := `
		INSERT INTO refresh_tokens (token, user_id,expires_at) VALUES ($1,$2,$3)
	`

	_, err := pool.Exec(context.TODO(), sql, token, id, time.Now().Add(60*24*time.Hour))
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetRefreshToken(pool *pgxpool.Pool, token string) (bool, uuid.UUID, error) { //todo giving no rows in result set for org
	var userId uuid.UUID
	err := pool.QueryRow(context.TODO(), "SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL", token).Scan(&userId)
	if err != nil {
		return false, uuid.UUID{}, err
	}
	if (userId == uuid.UUID{}) {
		return false, uuid.UUID{}, errors.New("user id is empty")
	}
	return true, userId, nil //Getting nil retu
}

func RevokeToken(pool *pgxpool.Pool, token string) error {
	_, err := pool.Exec(context.TODO(), "UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW() WHERE token = $1", token)
	if err != nil {
		return err
	}
	return nil
}
