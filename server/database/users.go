package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id            uuid.UUID
	Username      string
	HashPassword  string
	Email         string
	Date_of_birth time.Time
}

type UpdatedUser struct {
	Username string `json:"Username"`
}

type NewUser struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	// Email         string    `json:"Email"`
	// Date_of_birth time.Time `json:"Date_of_birth"`
}

// id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// username VARCHAR(40) NOT NULL,
// email VARCHAR(100) ,
// password VARCHAR(50) NOT NULL,
// date_of_birth DATE,
// country_id INT REFERENCES locations(id) ON DELETE SET NULL,
// state_id INT REFERENCES locations(id) ON DELETE SET NULL,
// city_id INT REFERENCES locations(id) ON DELETE SET NULL,
// profile_pi

func GetUser(username any) (User, error) {
	ctx, pool := createConnection()
	var user User
	err := pool.QueryRow(ctx, "SELECT id,username,password from users WHERE username = $1", username).Scan( // add email,date_of_birth
		&user.Id, &user.Username, &user.HashPassword,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed in getUser!: %v\n", err)
		return User{}, err
	}

	return user, nil
}

func GetUserById(id uuid.UUID) (User, error) {
	ctx, pool := createConnection()
	var user User
	err := pool.QueryRow(ctx, "SELECT id,username FROM users WHERE id = $1", id).Scan(
		&user.Id, &user.Username,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query row failed in getuserbyid %v\n", err)
		return User{}, err
	}
	return user, nil
}

func CreateUser(user NewUser, hashedPassword string) {
	ctx, pool := createConnection()

	exists, err := UserExists(ctx, pool, user.Username)
	if err != nil { //todo beef up this error handler
		fmt.Printf("error checking user exists: %v\n", err)
		return
	}
	if exists {
		return
	}
	sql := `INSERT INTO users(username,password) VALUES ($1,$2)`

	_, err = pool.Exec(ctx, sql, user.Username, hashedPassword) //add other options for new user like dob and email
	inserQueryFail(err, "creating user")
}

func UserExists(ctx context.Context, pool *pgxpool.Pool, username string) (bool, error) {
	rows, err := pool.Query(ctx, "SELECT * FROM users WHERE username = $1 LIMIT 1", username)
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

func UpdateUser(userId uuid.UUID, userInfo UpdatedUser) error {
	ctx, pool := createConnection()

	// check if user exists?
	sql := `UPDATE USERS SET username = $1 WHERE id = $2`

	_, err := pool.Exec(ctx, sql, userInfo.Username, userId)
	if err != nil {
		return err
	}
	return nil
}
