package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id              uuid.UUID
	Username        string
	HashPassword    string
	Email           string
	Profile_picture *string
	Date_of_birth   time.Time
	Is_organization bool
}

type NewUser struct {
	Id              uuid.UUID
	Username        string `json:"Username"`
	Password        string `json:"Password"`
	Email           string `json:"Email"`
	Is_organization bool   `json:"Is_organization"`
	// Date_of_birth time.Time `json:"Date_of_birth"`
}

func GetUser(pool *pgxpool.Pool, username any) (User, error) {
	var user User
	err := pool.QueryRow(context.TODO(), "SELECT id,username,password from users WHERE username = $1", username).Scan( // add email,date_of_birth
		&user.Id, &user.Username, &user.HashPassword,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed in getUser!: %v\n", err)
		return User{}, err
	}

	return user, nil
}

func GetUserById(pool *pgxpool.Pool, id uuid.UUID) (User, error) {
	var user User
	var dateOfBirth pgtype.Date

	rows := pool.QueryRow(context.TODO(), "SELECT id,username,email, profile_picture_src, date_of_birth FROM users WHERE id = $1", id)
	err := rows.Scan(
		&user.Id, &user.Username, &user.Email, &user.Profile_picture, &dateOfBirth,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query row failed in GetUserById %v\n", err)
		return User{}, err
	}

	if dateOfBirth.Status == pgtype.Present {
		user.Date_of_birth = dateOfBirth.Time
	} else {
		user.Date_of_birth = time.Now()
	}

	return user, nil
}

func CreateUser(pool *pgxpool.Pool, user User) (uuid.UUID, error) {

	exists, err := UserExists(pool, user.Username)
	if err != nil { //todo beef up this error handler
		fmt.Printf("error checking user exists: %v\n", err)
		return uuid.UUID{}, err
	}
	if exists {
		return uuid.UUID{}, errors.New("user already exists")
	}
	sql := `INSERT INTO users(username,password,email,profile_picture_src, is_organization) VALUES ($1,$2, $3, $4, $5) RETURNING id`

	var storedId uuid.UUID
	pool.QueryRow(context.TODO(), sql, user.Username, user.HashPassword, user.Email, user.Profile_picture, user.Is_organization).Scan(&storedId) //add other options for new user like dob and email
	if (storedId == uuid.UUID{}) {
		return uuid.UUID{}, errors.New("users Stored Id came out as empty")
	}
	return storedId, nil
}

func UserExists(pool *pgxpool.Pool, username string) (bool, error) {
	rows, err := pool.Query(context.TODO(), "SELECT * FROM users WHERE username = $1 LIMIT 1", username)
	if err != nil {
		return false, err
	}

	rowsProcessed := 0
	for rows.Next() {
		rowsProcessed++
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	if rowsProcessed < 1 {
		return false, nil
	}
	return true, nil
}

func UpdateUser(pool *pgxpool.Pool, userInfo User) error {

	// check if user exists?
	sql := `UPDATE users SET username = $1, email = $2, date_of_birth = $3, profile_picture_src = $4, updated_at = NOW() WHERE id = $5`

	_, err := pool.Exec(context.TODO(), sql, userInfo.Username, userInfo.Email, userInfo.Date_of_birth, userInfo.Profile_picture, userInfo.Id)
	if err != nil {
		return err
	}
	return nil
}
