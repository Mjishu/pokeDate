package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateRefreshToken(token string, user_id uuid.UUID) (bool, error) {
	ctx, pool := createConnection()

	sql := `
		INSERT INTO refresh_tokens (token, user_id,expires_at) VALUES ($1,$2,$3)
	`

	_, err := pool.Exec(ctx, sql, token, user_id, time.Now().Add(60*24*time.Hour))
	if err != nil {
		return false, errors.New("issue inserting refresh token into table")
	}
	return true, nil
}

func GetRefreshToken(token string) (bool, uuid.UUID) {
	ctx, pool := createConnection()
	var userId uuid.UUID
	err := pool.QueryRow(ctx, "SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()", token).Scan(&userId)
	if err != nil {
		fmt.Printf("error querying refreshToken %v\n", err)
		return false, uuid.UUID{}
	}
	if (userId == uuid.UUID{}) {
		return false, uuid.UUID{}
	}
	return true, userId //Getting nil retu
}
