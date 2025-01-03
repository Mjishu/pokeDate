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

func GetRefreshToken(token string) (bool, uuid.UUID) { //todo CHECK IF REVOKED_AT IS THERE
	ctx, pool := createConnection()
	var userId uuid.UUID
	err := pool.QueryRow(ctx, "SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL", token).Scan(&userId)
	if err != nil {
		fmt.Printf("error querying refreshToken %v\n", err)
		return false, uuid.UUID{}
	}
	if (userId == uuid.UUID{}) {
		return false, uuid.UUID{}
	}
	return true, userId //Getting nil retu
}

func RevokeToken(token string) error {
	ctx, pool := createConnection()
	_, err := pool.Exec(ctx, "UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW() WHERE token = $1", token)
	if err != nil {
		return errors.New("issue revoking token")
	}
	return nil
}
